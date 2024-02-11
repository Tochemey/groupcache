/*
 * Copyright 2024 Tochemey
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cluster defines the groupcache cluster mode
package cluster

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	goset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
	"github.com/tochemey/groupcache/v2"
	"github.com/tochemey/groupcache/v2/consistenthash"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/tochemey/groupcache/v2/log"
	"go.uber.org/atomic"
)

// Interface defines the Node interface
type Interface interface {
	// Start starts the Node engine
	Start(ctx context.Context) error
	// Stop stops the Node engine
	Stop(ctx context.Context) error
}

type Node struct {
	// specifies the discovery provider
	discoveryProvider discovery.Provider
	// specifies the discovery options
	discoveryOptions discovery.Config
	// specifies the logger
	logger log.Logger

	// specifies the cache pool
	pool *groupcache.HTTPPool

	// specifies the http server
	server *http.Server

	// specifies the Node self
	self *discovery.Node

	mu *sync.RWMutex
	// specifies the list of peers
	peers goset.Set[string]

	// specifies the hasher
	hasher consistenthash.Hash

	// specifies the replica count
	replicaCount int

	// Specifies the shutdown timeout. The default value is 30s
	shutdownTimeout time.Duration

	started *atomic.Bool
}

// enforce compilation error
var _ Interface = &Node{}

// New create an instance of the cluster node
func New(ctx context.Context, sd *discovery.ServiceDiscovery, opts ...Option) (*Node, error) {
	// create an instance of the Node
	node := &Node{
		discoveryProvider: sd.Provider(),
		discoveryOptions:  sd.Config(),
		logger:            log.DefaultLogger,
		mu:                &sync.RWMutex{},
		peers:             goset.NewSet[string](),
		started:           atomic.NewBool(false),
	}
	// apply the various options
	for _, opt := range opts {
		opt.Apply(node)
	}

	// get the host info
	hostNode, err := discovery.Self()
	// handle the error
	if err != nil {
		node.logger.Error(errors.Wrap(err, "failed get the host node.💥"))
		return nil, err
	}

	// set the host startClusterNode
	node.self = hostNode
	// add the host to the list of peers
	node.peers.Add(hostNode.PeerURL())

	return node, nil
}

// Start starts the groupcache cluster node
func (x *Node) Start(ctx context.Context) error {
	// acquire the lock
	x.mu.Lock()
	// release the lock when done
	defer x.mu.Unlock()

	// set the logger
	logger := x.logger

	// add some logging information
	logger.Infof("Starting groupcache cluster node on (%s)....🤔", x.self.Address())

	// no-op when already started
	if x.started.Load() {
		return nil
	}

	// set the service discovery config
	if err := x.discoveryProvider.SetConfig(x.discoveryOptions); err != nil {
		// log the error
		logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("unable to set the configuration")

		// return the error
		return errors.Wrapf(err, "failed to set config for the (%s) service discovery provider", x.discoveryProvider.ID())
	}

	// initialize the service discovery
	if err := x.discoveryProvider.Initialize(); err != nil {
		// log the error
		logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("failed to initialize")

		// return the error
		return errors.Wrapf(err, "failed to initialize the (%s) service discovery provider", x.discoveryProvider.ID())
	}

	// set the service discovery config
	if err := x.discoveryProvider.Register(); err != nil {
		// log the error
		logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("failed to register")

		// return the error
		return errors.Wrapf(err, "failed to register for the (%s) service discovery provider", x.discoveryProvider.ID())
	}

	// build the list of addresses
	peerURLs := goset.NewSet[string]()
	// let us perform some discovery
	nodes, err := x.discoveryProvider.DiscoverNodes()
	// handle the error
	if err != nil {
		// log the error
		logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("failed to discover initial nodes")

		return errors.Wrapf(err, "(%s) service discovery provider failed to find some nodes", x.discoveryProvider.ID())
	}

	// build the addresses list
	for _, node := range nodes {
		peerURLs.Add(node.PeerURL())
	}

	// add some debug logger
	x.logger.Debugf("peers discovered [%s]", strings.Join(peerURLs.ToSlice(), ","))

	// append the list of peers
	x.peers.Append(peerURLs.ToSlice()...)

	// create an instance of the cache pool
	x.pool = groupcache.NewHTTPPoolOpts(x.self.PeerURL(), &groupcache.HTTPPoolOptions{
		Replicas: x.replicaCount,
		HashFn:   x.hasher,
		Context: func(request *http.Request) context.Context {
			return ctx
		},
	})

	// add some debug logger
	x.logger.Debugf("node peers pool [%s]", strings.Join(x.peers.ToSlice(), ","))

	// attempt to connect to an existing cluster
	x.pool.Set(x.peers.ToSlice()...)

	// set the http server
	x.server = &http.Server{
		Addr: x.self.Address(),
		// The maximum duration for reading the entire request, including the body.
		// It’s implemented in net/http by calling SetReadDeadline immediately after Accept
		// ReadTimeout := handler_timeout + ReadHeaderTimeout + wiggle_room
		ReadTimeout: 3 * time.Second,
		// ReadHeaderTimeout is the amount of time allowed to read request headers
		ReadHeaderTimeout: time.Second,
		// WriteTimeout is the maximum duration before timing out writes of the response.
		// It is reset whenever a new request’s header is read.
		// This effectively covers the lifetime of the ServeHTTP handler stack
		WriteTimeout: time.Second,
		// IdleTimeout is the maximum amount of time to wait for the next request when keep-alive are enabled.
		// If IdleTimeout is zero, the value of ReadTimeout is used. Not relevant to request timeouts
		IdleTimeout: 1200 * time.Second,
		Handler:     x.pool,
		// Set the base context to incoming context of the system
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	// listen and service requests
	go func() {
		if err := x.server.ListenAndServe(); err != nil {
			// check the error type
			if !errors.Is(err, http.ErrServerClosed) {
				x.logger.Fatal(errors.Wrap(err, "failed to start server"))
			}
		}
	}()

	// read member events
	memberEvents, err := x.discoveryProvider.Watch(ctx)
	if err != nil {
		// TODO add some log
		return err
	}

	// handle member change events
	go x.handleDiscoveryEvents(memberEvents)

	// set the started
	x.started = atomic.NewBool(true)
	return nil
}

// Stop shutdown the cluster node
func (x *Node) Stop(ctx context.Context) error {
	// acquire the lock
	x.mu.Lock()
	// release the lock when done
	defer x.mu.Unlock()

	// no op when the cluster node has not started
	if !x.started.Load() {
		// add some warning log
		x.logger.Warn("groupcache cluster node not started")
		return nil
	}

	// create a cancellation context to gracefully shutdown
	ctx, cancel := context.WithTimeout(ctx, x.shutdownTimeout)
	defer cancel()

	// close the discovery provider
	if err := x.discoveryProvider.Close(); err != nil {
		x.logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("failed to close")

		return errors.Wrapf(err, "failed to close the (%s) service discovery provider", x.discoveryProvider.ID())
	}

	// stop the discovery provider
	if err := x.discoveryProvider.Deregister(); err != nil {
		x.logger.
			With("discovery", x.discoveryProvider.ID()).
			Error("failed to de-register")

		return errors.Wrapf(err, "failed to de-register the (%s) service discovery provider", x.discoveryProvider.ID())
	}

	// stop the server
	if err := x.server.Shutdown(ctx); err != nil {
		x.logger.Error("unable to stop peering server")
		return err
	}

	// set started to false
	x.started = atomic.NewBool(false)
	// reset the peers
	x.peers.Clear()
	return nil
}

// handleDiscoveryEvents handles discovery events
func (x *Node) handleDiscoveryEvents(events <-chan discovery.Event) {
	for event := range events {
		switch t := event.(type) {
		case *discovery.NodeAdded:
			// grab the added node
			node := t.Node
			// acquire the lock
			x.mu.Lock()
			// add the node to the list of peers
			x.peers.Add(node.Address())
			// update the pool set
			x.pool.Set(x.peers.ToSlice()...)
			// release the lock
			x.mu.Unlock()
		case *discovery.NodeModified:
			// grab the updated node
			node := t.Node
			// grab the old node
			current := t.Current
			// acquire the lock
			x.mu.Lock()
			// remove the list of peers
			x.peers.Remove(current.Address())
			// add the node to the list of peers
			x.peers.Add(node.Address())
			// update the pool set
			x.pool.Set(x.peers.ToSlice()...)
			// release the lock
			x.mu.Unlock()
		case *discovery.NodeRemoved:
			// grab the removed node
			node := t.Node
			// acquire the lock
			x.mu.Lock()
			// remove the list of peers
			x.peers.Remove(node.Address())
			// update the pool set
			x.pool.Set(x.peers.ToSlice()...)
			// release the lock
			x.mu.Unlock()
		}
	}
}
