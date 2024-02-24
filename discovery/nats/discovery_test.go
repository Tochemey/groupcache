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
	"os"
	"strconv"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/tochemey/groupcache/v2/log"
	"github.com/travisjeffery/go-dynaport"
	"go.uber.org/atomic"
)

func startNatsServer(t *testing.T) *natsserver.Server {
	t.Helper()
	serv, err := natsserver.NewServer(&natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	})

	require.NoError(t, err)

	ready := make(chan bool)
	go func() {
		ready <- true
		serv.Start()
	}()
	<-ready

	if !serv.ReadyForConnections(2 * time.Second) {
		t.Fatalf("nats-io server failed to start")
	}

	return serv
}

func newPeer(t *testing.T, serverAddr string) *Discovery {
	// generate the ports for the single node
	nodePorts := dynaport.Get(3)
	groupPort := nodePorts[0]

	// create a Cluster node
	host := "localhost"
	// set the environments
	require.NoError(t, os.Setenv("GROUP_PORT", strconv.Itoa(groupPort)))
	require.NoError(t, os.Setenv("NODE_NAME", "testNode"))
	require.NoError(t, os.Setenv("NODE_IP", host))

	// create the various config option
	applicationName := "test"
	natsSubject := "some-subject"
	// create the instance of provider
	provider := NewDiscovery()

	// create the config
	config := discovery.Options{
		ApplicationName: applicationName,
		Server:          serverAddr,
		Subject:         natsSubject,
	}

	// set config
	err := provider.SetConfig(config)
	require.NoError(t, err)

	// initialize
	err = provider.Initialize()
	require.NoError(t, err)
	// clear the env var
	require.NoError(t, os.Unsetenv("GROUP_PORT"))
	require.NoError(t, os.Unsetenv("NODE_NAME"))
	require.NoError(t, os.Unsetenv("NODE_IP"))
	// return the provider
	return provider
}

func TestDiscovery(t *testing.T) {
	t.Run("With a new instance", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		require.NotNil(t, provider)
		// assert that provider implements the Discovery interface
		// this is a cheap test
		// assert the type of svc
		assert.IsType(t, &Discovery{}, provider)
		var p interface{} = provider
		_, ok := p.(discovery.Provider)
		assert.True(t, ok)
	})
	t.Run("With ID assertion", func(t *testing.T) {
		// cheap test
		// create the instance of provider
		provider := NewDiscovery()
		require.NotNil(t, provider)
		assert.Equal(t, "nats", provider.ID())
	})
	t.Run("With SetConfig", func(t *testing.T) {
		// create the various config option
		natsServer := "nats://127.0.0.1:2322"
		applicationName := "test"
		natsSubject := "some-subject"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Options{
			ApplicationName: applicationName,
			Server:          natsServer,
			Subject:         natsSubject,
		}

		// set config
		assert.NoError(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: already initialized", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)
		// create the various config option
		natsServer := srv.Addr().String()
		applicationName := "tests"
		natsSubject := "some-subject"
		// create the instance of provider
		provider := NewDiscovery()
		provider.initialized = atomic.NewBool(true)
		// create the config
		config := discovery.Options{
			ApplicationName: applicationName,
			Server:          natsServer,
			Subject:         natsSubject,
		}

		// set config
		err := provider.SetConfig(config)
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrAlreadyInitialized.Error())
		// stop the NATS server
		t.Cleanup(srv.Shutdown)
	})
	t.Run("With SetConfig: nats server not set", func(t *testing.T) {
		applicationName := "tests"
		natsSubject := "some-subject"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Options{
			ApplicationName: applicationName,
			Subject:         natsSubject,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: nats subject not set", func(t *testing.T) {
		// create the various config option
		natsServer := "nats://127.0.0.1:2322"
		applicationName := "accounts"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Options{
			ApplicationName: applicationName,
			Server:          natsServer,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: application name not set", func(t *testing.T) {
		// create the various config option
		natsServer := "nats://127.0.0.1:2322"

		natsSubject := "some-subject"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Options{
			Server:  natsServer,
			Subject: natsSubject,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With Initialize", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(3)
		groupPort := nodePorts[0]

		// create a Cluster node
		host := "localhost"
		// set the environments
		t.Setenv("GROUP_PORT", strconv.Itoa(groupPort))
		t.Setenv("NODE_NAME", "testNode")
		t.Setenv("NODE_IP", host)

		// create the various config option
		natsServer := srv.Addr().String()
		applicationName := "test"
		natsSubject := "some-subject"
		// create the instance of provider
		provider := NewDiscovery(WithLogger(log.DiscardLogger))

		// create the config
		config := discovery.Options{
			ApplicationName: applicationName,
			Server:          natsServer,
			Subject:         natsSubject,
		}

		// set config
		err := provider.SetConfig(config)
		require.NoError(t, err)

		// initialize
		err = provider.Initialize()
		assert.NoError(t, err)

		// stop the NATS server
		t.Cleanup(srv.Shutdown)
	})
	t.Run("With Initialize: already initialized", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		provider.initialized = atomic.NewBool(true)
		assert.Error(t, provider.Initialize())
	})
	t.Run("With Initialize: with host config not set", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		assert.Error(t, provider.Initialize())
	})
	t.Run("With Register: already registered", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		provider.registered = atomic.NewBool(true)
		err := provider.Register()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrAlreadyRegistered.Error())
	})
	t.Run("With Deregister: already not registered", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		err := provider.Deregister()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
	})
	t.Run("With DiscoverNodes", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)
		// create two peers
		client1 := newPeer(t, srv.Addr().String())
		client2 := newPeer(t, srv.Addr().String())

		// no discovery is allowed unless registered
		peers, err := client1.DiscoverNodes()
		require.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
		require.Empty(t, peers)

		peers, err = client2.DiscoverNodes()
		require.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
		require.Empty(t, peers)

		// register client 2
		require.NoError(t, client2.Register())
		peers, err = client2.DiscoverNodes()
		require.NoError(t, err)
		require.Empty(t, peers)

		// register client 1
		require.NoError(t, client1.Register())
		peers, err = client1.DiscoverNodes()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.Len(t, peers, 1)
		discoveredNodeAddr := client2.hostNode.Address()
		require.Equal(t, peers[0].Address(), discoveredNodeAddr)

		// discover more peers from client 2
		peers, err = client2.DiscoverNodes()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.Len(t, peers, 1)
		discoveredNodeAddr = client1.hostNode.Address()
		require.Equal(t, peers[0].Address(), discoveredNodeAddr)

		// de-register client 2 but it can see client1
		require.NoError(t, client2.Deregister())
		peers, err = client2.DiscoverNodes()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		discoveredNodeAddr = client1.hostNode.Address()
		require.Equal(t, peers[0].Address(), discoveredNodeAddr)

		// client-1 cannot see the deregistered client2
		peers, err = client1.DiscoverNodes()
		require.NoError(t, err)
		require.Empty(t, peers)

		require.NoError(t, client1.Close())
		require.NoError(t, client2.Close())

		// stop the NATS server
		t.Cleanup(srv.Shutdown)
	})
	t.Run("With DiscoverPeers: not initialized", func(t *testing.T) {
		provider := NewDiscovery()
		peers, err := provider.DiscoverNodes()
		assert.Error(t, err)
		assert.Empty(t, peers)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
}
