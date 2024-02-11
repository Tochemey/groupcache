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

package kubernetes

import (
	"sort"
	"testing"
	"time"

	goset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tochemey/groupcache/v2/discovery"
	"go.uber.org/atomic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
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
		assert.Equal(t, "kubernetes", provider.ID())
	})
	t.Run("With DiscoverNodes", func(t *testing.T) {
		// create the namespace
		ns := "test"
		appName := "test"
		ts1 := time.Now()
		ts2 := time.Now()

		// create some bunch of mock pods
		pods := []runtime.Object{
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod1",
					Namespace: ns,
					Labels: map[string]string{
						"app.kubernetes.io/part-of":   appName,
						"app.kubernetes.io/component": appName,
						"app.kubernetes.io/name":      appName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Ports: []corev1.ContainerPort{
								{
									Name:          GroupCachePortName,
									ContainerPort: 3379,
								},
							},
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					PodIP: "10.0.0.23",
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.PodReady,
							Status: corev1.ConditionTrue,
						},
					},
					StartTime: &metav1.Time{
						Time: ts1,
					},
				},
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod2",
					Namespace: ns,
					Labels: map[string]string{
						"app.kubernetes.io/part-of":   appName,
						"app.kubernetes.io/component": appName,
						"app.kubernetes.io/name":      appName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Ports: []corev1.ContainerPort{
								{
									Name:          GroupCachePortName,
									ContainerPort: 3379,
								},
							},
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					PodIP: "10.0.0.24",
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.PodReady,
							Status: corev1.ConditionTrue,
						},
					},
					StartTime: &metav1.Time{
						Time: ts2,
					},
				},
			},
		}
		// create a mock kubernetes client
		client := testclient.NewSimpleClientset(pods...)
		// create the kubernetes discovery provider
		provider := Discovery{
			client:      client,
			initialized: atomic.NewBool(true),
			option: &option{
				NameSpace:       ns,
				ApplicationName: appName,
			},
		}
		// discover some nodes
		actual, err := provider.DiscoverNodes()
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.NotEmpty(t, actual)
		require.Len(t, actual, 2)

		expected := []string{
			"10.0.0.23:3379",
			"10.0.0.24:3379",
		}

		actualAddresses := goset.NewSet[string]()
		for _, node := range actual {
			actualAddresses.Add(node.Address())
		}
		addresses := actualAddresses.ToSlice()
		sort.Strings(addresses)
		sort.Strings(expected)

		assert.ElementsMatch(t, expected, addresses)
		assert.NoError(t, provider.Close())
	})
	t.Run("With DiscoverNodes: not initialized", func(t *testing.T) {
		provider := NewDiscovery()
		peers, err := provider.DiscoverNodes()
		assert.Error(t, err)
		assert.Empty(t, peers)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
	t.Run("With SetConfig", func(t *testing.T) {
		// create the various config option
		namespace := "default"
		applicationName := "accounts"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			ApplicationName: applicationName,
			Namespace:       namespace,
		}

		// set config
		assert.NoError(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: already initialized", func(t *testing.T) {
		// create the various config option
		namespace := "default"
		applicationName := "accounts"
		// create the instance of provider
		provider := NewDiscovery()
		provider.initialized = atomic.NewBool(true)
		// create the config
		config := discovery.Config{
			ApplicationName: applicationName,
			Namespace:       namespace,
		}

		// set config
		err := provider.SetConfig(config)
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrAlreadyInitialized.Error())
	})

	t.Run("With SetConfig: application name not set", func(t *testing.T) {
		// create the various config option
		namespace := "default"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			Namespace: namespace,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With SetConfig: namespace not set", func(t *testing.T) {
		applicationName := "accounts"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			ApplicationName: applicationName,
		}

		// set config
		assert.Error(t, provider.SetConfig(config))
	})
	t.Run("With Initialize", func(t *testing.T) {
		// create the various config option
		namespace := "default"
		applicationName := "accounts"
		// create the instance of provider
		provider := NewDiscovery()
		// create the config
		config := discovery.Config{
			ApplicationName: applicationName,
			Namespace:       namespace,
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
}
