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

package mdns

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/travisjeffery/go-dynaport"
	"go.uber.org/atomic"
)

func TestDiscovery(t *testing.T) {
	t.Run("With new instance", func(t *testing.T) {
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
		assert.Equal(t, "mdns", provider.ID())
	})
	t.Run("With SetConfig", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "AccountsSystem"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}

		// set config
		assert.NoError(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig with service name not set", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service: serviceType,
			Domain:  domain,
			Port:    port,
			IPv6:    false,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig with service not set", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceName := "AccountsSystem"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig with port not set", func(t *testing.T) {
		// create the various config option
		serviceType := "_workstation._tcp"
		serviceName := "AccountsSystem"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			IPv6:        false,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig with domain not set", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "AccountsSystem"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Port:        port,
			IPv6:        false,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig with ipv6 not set", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "AccountsSystem"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: already initialized", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "test"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		provider.initialized = atomic.NewBool(true)
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}

		// set config
		err := provider.SetConfig(config)
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrAlreadyInitialized.Error())
	})
	t.Run("With Initialize", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "test"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}
		// set config
		assert.NoError(t, provider.SetConfig(config))
		assert.NoError(t, provider.Initialize())
	})
	t.Run("With Initialize: already initialized", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		provider.initialized = atomic.NewBool(true)
		assert.Error(t, provider.Initialize())
	})
	t.Run("With Register", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "test"
		domain := "local"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}
		require.NoError(t, provider.SetConfig(config))
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())
		require.NoError(t, provider.Deregister())
		time.Sleep(time.Second)
		assert.False(t, provider.initialized.Load())
	})
	t.Run("With Register when already registered", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := strconv.Itoa(ports[0])
		serviceType := "_workstation._tcp"
		serviceName := "test"
		domain := "local"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     serviceType,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}
		require.NoError(t, provider.SetConfig(config))
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())
		err := provider.Register()
		require.Error(t, err)
		require.EqualError(t, err, discovery.ErrAlreadyRegistered.Error())
		require.NoError(t, provider.Deregister())
		time.Sleep(time.Second)
		assert.False(t, provider.initialized.Load())
	})
	t.Run("With Deregister", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		// for the sake of the test
		provider.initialized = atomic.NewBool(true)
		assert.NoError(t, provider.Deregister())
	})
	t.Run("With Deregister when not initialized", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery()
		// for the sake of the test
		provider.initialized = atomic.NewBool(false)
		err := provider.Deregister()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
	t.Run("With DiscoverPeers", func(t *testing.T) {
		// create the various config option
		ports := dynaport.Get(1)
		port := ports[0]
		service := "_workstation._tcp"
		serviceName := "test"
		domain := "local."
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Service:     service,
			ServiceName: serviceName,
			Domain:      domain,
			Port:        port,
			IPv6:        false,
		}
		require.NoError(t, provider.SetConfig(config))
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		// wait for registration to be completed
		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())

		// discover peers
		peers, err := provider.DiscoverNodes()
		require.NoError(t, err)
		require.NotEmpty(t, peers)

		assert.NoError(t, provider.Deregister())
	})
	t.Run("With DiscoverPeers: not initialized", func(t *testing.T) {
		provider := NewDiscovery()
		peers, err := provider.DiscoverNodes()
		assert.Error(t, err)
		assert.Empty(t, peers)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
}
