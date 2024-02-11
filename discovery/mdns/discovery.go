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

// Package mdns defines the mdns discovery provider
package mdns

import (
	"context"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
	"github.com/tochemey/groupcache/v2/discovery"
	"go.uber.org/atomic"
)

const (
	ServiceName = "name"
	Service     = "service"
	Domain      = "domain"
	Port        = "group-port"
	IPv6        = "ipv6"
)

// Option represents the mDNS provider option
type Option struct {
	// Provider specifies the provider name
	Provider string
	// Service specifies the service name
	ServiceName string
	// Service specifies the service type
	Service string
	// Specifies the service domain
	Domain string
	// Port specifies the port the service is listening to
	Port int
	// IPv6 states whether to fetch ipv6 address instead of ipv4
	IPv6 *bool
}

// Discovery defines the mDNS discovery provider
type Discovery struct {
	option *Option
	mu     sync.Mutex

	stopChan    chan struct{}
	publicChan  chan discovery.Event
	initialized *atomic.Bool

	// resolver is used to browse for service discovery
	resolver *zeroconf.Resolver

	server *zeroconf.Server
}

// enforce compilation error
var _ discovery.Provider = &Discovery{}

// NewDiscovery returns an instance of the mDNS discovery provider
func NewDiscovery() *Discovery {
	// create an instance of
	d := &Discovery{
		mu:          sync.Mutex{},
		stopChan:    make(chan struct{}, 1),
		publicChan:  make(chan discovery.Event, 2),
		initialized: atomic.NewBool(false),
		option:      &Option{},
	}

	return d
}

// ID returns the discovery provider identifier
func (d *Discovery) ID() string {
	return "mdns"
}

// Initialize the discovery provider
func (d *Discovery) Initialize() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()
	// first check whether the discovery provider is running
	if d.initialized.Load() {
		return discovery.ErrAlreadyInitialized
	}

	// check the options
	if d.option.Provider == "" {
		d.option.Provider = d.ID()
	}

	return nil
}

// Register registers this node to a service discovery directory.
func (d *Discovery) Register() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider has started
	// avoid to re-register the discovery
	if d.initialized.Load() {
		return discovery.ErrAlreadyRegistered
	}

	// initialize the resolver
	res, err := zeroconf.NewResolver(nil)
	// handle the error
	if err != nil {
		return errors.Wrap(err, "failed to instantiate the mDNS discovery provider")
	}
	// set the resolver
	d.resolver = res

	// register the service
	srv, err := zeroconf.Register(d.option.ServiceName, d.option.Service, d.option.Domain, d.option.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	// handle the error
	if err != nil {
		return err
	}

	// set the server
	d.server = srv

	// set initialized
	d.initialized = atomic.NewBool(true)
	return nil
}

// Deregister removes this node from a service discovery directory.
func (d *Discovery) Deregister() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider has started
	if !d.initialized.Load() {
		return discovery.ErrNotInitialized
	}

	// shutdown the registered service
	if d.server != nil {
		d.server.Shutdown()
	}

	// set the initialized to false
	d.initialized = atomic.NewBool(false)
	// stop the watchers
	close(d.stopChan)
	// close the public channel
	close(d.publicChan)
	// return
	return nil
}

// Close closes the provider
func (d *Discovery) Close() error {
	return nil
}

// SetConfig registers the underlying discovery configuration
func (d *Discovery) SetConfig(config discovery.Config) error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider is running
	if d.initialized.Load() {
		return discovery.ErrAlreadyInitialized
	}

	var err error
	// validate the meta
	// let us make sure we have the required options set

	// assert the presence of service instance
	if _, ok := config[ServiceName]; !ok {
		return errors.New("mDNS service name is not provided")
	}

	// assert the presence of the service type
	if _, ok := config[Service]; !ok {
		return errors.New("mDNS service type is not provided")
	}

	// assert the presence of the listening port
	if _, ok := config[Port]; !ok {
		return errors.New("mDNS listening port is not provided")
	}

	// assert the service domain
	if _, ok := config[Domain]; !ok {
		return errors.New("mDNS domain is not provided")
	}

	// assert the ipv6 domain
	if _, ok := config[IPv6]; !ok {
		return errors.New("mDNS ipv6 option is not provided")
	}

	// set the options
	if err = d.setOptions(config); err != nil {
		return errors.Wrap(err, "failed to instantiate the mDNS discovery provider")
	}
	return nil
}

// Watch returns event based upon node lifecycle
func (d *Discovery) Watch(ctx context.Context) (<-chan discovery.Event, error) {
	// first check whether the discovery provider is running
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}
	// run the watcher
	go d.watchPods(ctx)
	return d.publicChan, nil
}

// DiscoverNodes returns a list of known nodes.
func (d *Discovery) DiscoverNodes() ([]*discovery.Node, error) {
	// first check whether the discovery provider is running
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}
	// init entries channel
	entries := make(chan *zeroconf.ServiceEntry, 100)

	// create a context to browse the services for 5 seconds
	// TODO: make the timeout configurable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// let us browse the services
	if err := d.resolver.Browse(ctx, d.option.Service, d.option.Domain, entries); err != nil {
		return nil, err
	}
	<-ctx.Done()

	// set ipv6 filter
	v6 := false
	if d.option.IPv6 != nil {
		v6 = *d.option.IPv6
	}

	// define the addresses list
	nodes := make([]*discovery.Node, 0, len(entries))
	for entry := range entries {
		// validate the entry
		if !d.validateEntry(entry) {
			continue
		}
		// lookup for v6 address
		if v6 {
			// iterate the list of ports
			for _, addr := range entry.AddrIPv6 {
				nodes = append(nodes, discovery.NewNode(entry.ServiceInstanceName(), addr.String(), uint32(entry.Port)))
			}
		}

		// iterate the list of ports
		for _, addr := range entry.AddrIPv4 {
			nodes = append(nodes, discovery.NewNode(entry.ServiceInstanceName(), addr.String(), uint32(entry.Port)))
		}
	}
	return nodes, nil
}

// setOptions sets the kubernetes discoConfig
func (d *Discovery) setOptions(config discovery.Config) (err error) {
	// create an instance of Option
	option := new(Option)
	// extract the service name
	option.ServiceName, err = config.GetString(ServiceName)
	// handle the error in case the service instance value is not properly set
	if err != nil {
		return err
	}
	// extract the service name
	option.Service, err = config.GetString(Service)
	// handle the error in case the service type value is not properly set
	if err != nil {
		return err
	}
	// extract the service domain
	option.Domain, err = config.GetString(Domain)
	// handle the error when the domain is not properly set
	if err != nil {
		return err
	}
	// extract the port the service is listening to
	option.Port, err = config.GetInt(Port)
	// handle the error when the port is not properly set
	if err != nil {
		return err
	}

	// extract the type of ip address to lookup
	option.IPv6, err = config.GetBool(IPv6)
	// handle the error
	if err != nil {
		return err
	}

	// in case none of the above extraction fails then set the option
	d.option = option
	return nil
}

// validateEntry validates the mDNS discovered entry
func (d *Discovery) validateEntry(entry *zeroconf.ServiceEntry) bool {
	return entry.Port == d.option.Port &&
		entry.Service == d.option.Service &&
		entry.Domain == d.option.Domain &&
		entry.Instance == d.option.ServiceName
}

// watchPods keeps a watch on mDNS pods activities and emit
// respective event when needed
func (d *Discovery) watchPods(ctx context.Context) {
	// set ipv6 filter
	v6 := false
	if d.option.IPv6 != nil {
		v6 = *d.option.IPv6
	}
	// create an array of zeroconf service entries
	entries := make(chan *zeroconf.ServiceEntry, 5)
	// create a go routine that will push discovery event to the discovery channel
	go func() {
		for {
			select {
			case <-d.stopChan:
				return
			case entry := <-entries:
				if v6 {
					discovery.NewNode(entry.ServiceInstanceName(), entry.AddrIPv6[0].String(), uint32(entry.Port))
					continue
				}
				node := discovery.NewNode(entry.ServiceInstanceName(), entry.AddrIPv4[0].String(), uint32(entry.Port))
				event := &discovery.NodeAdded{Node: node}
				// add to the channel
				d.publicChan <- event
			}
		}
	}()

	// wrap the context in a cancellation context
	ctx, cancel := context.WithCancel(ctx)
	for {
		select {
		case <-d.stopChan:
			cancel()
			return
		default:
			// browse for all services of a given type in a given domain.
			if err := d.resolver.Browse(ctx, d.option.Service, d.option.Domain, entries); err != nil {
				// log the error
				panic(errors.Wrap(err, "failed to lookup"))
			}
		}
	}
}
