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

package nats

import (
	"context"
	"sync"
	"time"

	"github.com/flowchartsman/retry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/tochemey/groupcache/v2/log"
	"go.uber.org/atomic"
)

const (
	Server          string = "nats-server"  // Server specifies the Nats server Address
	Subject         string = "nats-subject" // Subject specifies the NATs subject
	ApplicationName        = "app_name"     // ApplicationName specifies the application name. This often matches the actor system name
	Timeout                = "timeout"      // Timeout specifies the discovery timeout. The default value is 1 second
)

// option represents the nats provider option
type option struct {
	// NatsServer defines the nats server
	// nats://host:port of a nats server
	NatsServer string
	// NatsSubject defines the custom NATS subject
	NatsSubject string
	// ApplicationName specifies the running application
	ApplicationName string
	// Timeout defines the nodes discovery timeout
	Timeout time.Duration
}

// Discovery represents the kubernetes discovery
type Discovery struct {
	config *option
	mu     sync.Mutex

	initialized *atomic.Bool
	registered  *atomic.Bool

	// define the nats connection
	natsConnection *nats.EncodedConn
	// define a slice of subscriptions
	subscriptions []*nats.Subscription

	// defines the host node
	hostNode *discovery.Node

	publicChan chan discovery.Event
	stopChan   chan struct{}
	// define a logger
	logger log.Logger
}

// enforce compilation error
var _ discovery.Provider = &Discovery{}

// NewDiscovery returns an instance of the kubernetes discovery provider
func NewDiscovery(opts ...Option) *Discovery {
	// create an instance of
	d := &Discovery{
		mu:          sync.Mutex{},
		initialized: atomic.NewBool(false),
		registered:  atomic.NewBool(false),
		config:      &option{},
		logger:      log.DefaultLogger,
		stopChan:    make(chan struct{}, 1),
		publicChan:  make(chan discovery.Event, 2),
	}

	// apply the various options
	for _, opt := range opts {
		opt.Apply(d)
	}

	return d
}

// ID returns the discovery provider id
func (d *Discovery) ID() string {
	return "nats"
}

// Initialize initializes the plugin: registers some internal data structures, clients etc.
func (d *Discovery) Initialize() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider is running
	if d.initialized.Load() {
		return discovery.ErrAlreadyInitialized
	}

	// set the default discovery timeout
	if d.config.Timeout <= 0 {
		d.config.Timeout = time.Second
	}

	// grab the host node
	hostNode, err := discovery.Self()
	// handle the error
	if err != nil {
		return err
	}

	// create the nats connection option
	opts := nats.GetDefaultOptions()
	opts.Url = d.config.NatsServer
	//opts.Servers = n.Config.Servers
	opts.Name = hostNode.Name()
	opts.ReconnectWait = 2 * time.Second
	opts.MaxReconnect = -1

	var (
		connection *nats.Conn
	)

	// let us connect using an exponential backoff mechanism
	// we will attempt to connect 5 times to see whether there have started successfully or not.
	const (
		maxRetries = 5
	)
	// create a new instance of retrier that will try a maximum of five times, with
	// an initial delay of 100 ms and a maximum delay of opts.ReconnectWait
	retrier := retry.NewRetrier(maxRetries, 100*time.Millisecond, opts.ReconnectWait)
	// handle the retry error
	err = retrier.Run(func() error {
		// connect to the NATs server
		connection, err = opts.Connect()
		if err != nil {
			return err
		}
		// successful connection
		return nil
	})

	// create the NATs connection encoder
	encodedConn, err := nats.NewEncodedConn(connection, protobuf.PROTOBUF_ENCODER)
	// handle the error
	if err != nil {
		return err
	}
	// set the connection
	d.natsConnection = encodedConn
	// set initialized
	d.initialized = atomic.NewBool(true)
	// set the host node
	d.hostNode = hostNode
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
	if d.registered.Load() {
		return discovery.ErrAlreadyRegistered
	}

	// get the host node gossip address
	me := d.hostNode.Address()

	// create the subscription handler
	subscriptionHandler := func(subj, reply string, msg *Message) {
		switch msg.GetMessageType() {
		case MessageType_MESSAGE_TYPE_DEREGISTER:
			// add logging information
			d.logger.Infof("received an de-registration request from peer[name=%s, host=%s, port=%d]",
				msg.GetName(), msg.GetHost(), msg.GetPort())

			// create the node found
			node := discovery.NewNode(msg.GetName(), msg.GetHost(), uint32(msg.GetPort()))

			// get the found peer address
			addr := node.Address()
			// ignore this node
			if addr != me {
				// here we find a node let us raise the node registered event
				event := &discovery.NodeRemoved{Node: node}
				// add to the channel
				d.publicChan <- event
			}

		case MessageType_MESSAGE_TYPE_REGISTER:
			// add logging information
			d.logger.Infof("received an registration request from peer[name=%s, host=%s, port=%d]",
				msg.GetName(), msg.GetHost(), msg.GetPort())

			// create the node found
			node := discovery.NewNode(msg.GetName(), msg.GetHost(), uint32(msg.GetPort()))

			// get the found peer address
			addr := node.Address()
			// ignore this node
			if addr != me {
				// here we find a node let us raise the node registered event
				event := &discovery.NodeAdded{Node: node}
				// add to the channel
				d.publicChan <- event
			}

		case MessageType_MESSAGE_TYPE_REQUEST:
			// create the node found
			node := discovery.NewNode(msg.GetName(), msg.GetHost(), uint32(msg.GetPort()))

			// exclude self
			if me != node.Address() {
				// add logging information
				d.logger.Infof("received an identification request from peer[name=%s, host=%s, port=%d]",
					msg.GetName(), msg.GetHost(), msg.GetPort())
				// send the reply
				replyMessage := &Message{
					Host:        d.hostNode.Host(),
					Port:        int32(d.hostNode.Port()),
					Name:        d.hostNode.Name(),
					MessageType: MessageType_MESSAGE_TYPE_RESPONSE,
				}

				// send the reply and handle the error
				if err := d.natsConnection.Publish(reply, replyMessage); err != nil {
					d.logger.Errorf("failed to reply for identification request from peer[name=%s, host=%s, port=%d]",
						msg.GetName(), msg.GetHost(), msg.GetPort())
				}
			}
		}
	}
	// start listening to incoming messages
	subscription, err := d.natsConnection.Subscribe(d.config.NatsSubject, subscriptionHandler)
	// return any eventual error
	if err != nil {
		return err
	}

	// add the subscription to the list of subscriptions
	d.subscriptions = append(d.subscriptions, subscription)
	// set the registration flag to true
	d.registered = atomic.NewBool(true)
	return nil
}

// Deregister removes this node from a service discovery directory.
func (d *Discovery) Deregister() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()
	// first check whether the discovery provider has been registered or not
	if !d.registered.Load() {
		return discovery.ErrNotRegistered
	}

	// shutdown all the subscriptions
	for _, subscription := range d.subscriptions {
		// when subscription is defined
		if subscription != nil {
			// check whether the subscription is active or not
			if subscription.IsValid() {
				// unsubscribe and return when there is an error
				if err := subscription.Unsubscribe(); err != nil {
					return err
				}
			}
		}
	}

	// send the de-registration message to notify peers
	if d.natsConnection != nil {
		// send a message to deregister stating we are out
		return d.natsConnection.Publish(d.config.NatsSubject, &Message{
			Host:        d.hostNode.Host(),
			Port:        int32(d.hostNode.Port()),
			Name:        d.hostNode.Name(),
			MessageType: MessageType_MESSAGE_TYPE_DEREGISTER,
		})
	}

	// set registered to false
	d.registered.Store(false)
	// stop the watchers
	close(d.stopChan)
	// close the public channel
	close(d.publicChan)

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

	return d.setConfig(config)
}

// DiscoverNodes returns a list of known nodes.
func (d *Discovery) DiscoverNodes() ([]*discovery.Node, error) {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider is running
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}

	// later check the provider is registered
	if !d.registered.Load() {
		return nil, discovery.ErrNotRegistered
	}

	// Set up a reply channel, then broadcast for all peers to
	// report their presence.
	// collect as many responses as possible in the given timeout.
	inbox := nats.NewInbox()
	recv := make(chan *Message)

	// bind to receive messages
	sub, err := d.natsConnection.BindRecvChan(inbox, recv)
	if err != nil {
		return nil, err
	}

	// send registration request and return in case of error
	if err = d.natsConnection.PublishRequest(d.config.NatsSubject, inbox, &Message{
		Host:        d.hostNode.Host(),
		Port:        int32(d.hostNode.Port()),
		Name:        d.hostNode.Name(),
		MessageType: MessageType_MESSAGE_TYPE_REQUEST,
	}); err != nil {
		return nil, err
	}

	// define the list of peers to lookup
	var peers []*discovery.Node
	// define the timeout
	timeout := time.After(d.config.Timeout)
	// get the host node gossip address
	me := d.hostNode.Address()

	// loop till we have enough nodes
	for {
		select {
		case m, ok := <-recv:
			if !ok {
				// Subscription is closed
				return peers, nil
			}

			// create the node found
			node := discovery.NewNode(m.GetName(), m.GetHost(), uint32(m.GetPort()))

			// get the found peer address
			addr := node.Address()
			// ignore this node
			if addr == me {
				// Ignore a reply from self
				continue
			}

			// add the peer to the list of discovered peers
			peers = append(peers, node)

		case <-timeout:
			// close the receiving channel
			_ = sub.Unsubscribe()
			close(recv)
		}
	}
}

// Close closes the provider
func (d *Discovery) Close() error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// set the initialized to false
	d.initialized.Store(false)

	if d.natsConnection != nil {
		defer func() {
			// close the underlying connection
			d.natsConnection.Close()
			d.natsConnection = nil
		}()

		// shutdown all the subscriptions
		for _, subscription := range d.subscriptions {
			// when subscription is defined
			if subscription != nil {
				// check whether the subscription is active or not
				if subscription.IsValid() {
					// unsubscribe and return when there is an error
					if err := subscription.Unsubscribe(); err != nil {
						return err
					}
				}
			}
		}
		// flush all messages
		if err := d.natsConnection.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Watch returns event based upon node lifecycle
func (d *Discovery) Watch(ctx context.Context) (<-chan discovery.Event, error) {
	// first check whether the discovery provider is running
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}
	return d.publicChan, nil
}

// setConfig sets the kubernetes option
func (d *Discovery) setConfig(config discovery.Config) (err error) {
	// create an instance of option
	option := new(option)
	// extract the nats server address
	option.NatsServer, err = config.GetString(Server)
	// handle the error in case the nats server value is not properly set
	if err != nil {
		return err
	}
	// extract the nats subject address
	option.NatsSubject, err = config.GetString(Subject)
	// handle the error in case the nats subject value is not properly set
	if err != nil {
		return err
	}
	// extract the application name
	option.ApplicationName, err = config.GetString(ApplicationName)
	// handle the error in case the application name value is not properly set
	if err != nil {
		return err
	}
	// in case none of the above extraction fails then set the option
	d.config = option
	return nil
}
