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

// Package kubernetes defines the kubernetes discovery provider
package kubernetes

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tochemey/groupcache/v2/discovery"
	"go.uber.org/atomic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/utils/strings/slices"
)

const (
	Namespace          string = "namespace"  // Namespace specifies the kubernetes namespace
	ApplicationName           = "app_name"   // ApplicationName specifies the application name.
	GroupCachePortName        = "cache-port" // GroupCachePortName specifies the group cache port name
)

// option defines the k8 discovery option
type option struct {
	// Provider specifies the provider name
	Provider string
	// NameSpace specifies the namespace
	NameSpace string
	// ApplicationName specifies the running application
	ApplicationName string
}

// Discovery represents the kubernetes discovery
type Discovery struct {
	option *option
	client kubernetes.Interface
	mu     sync.Mutex

	stopChan   chan struct{}
	publicChan chan discovery.Event

	// states whether discovery provider has started or not
	initialized *atomic.Bool
}

// enforce compilation error
var _ discovery.Provider = &Discovery{}

// NewDiscovery returns an instance of the kubernetes discovery provider
func NewDiscovery() *Discovery {
	// create an instance of
	discovery := &Discovery{
		mu:          sync.Mutex{},
		stopChan:    make(chan struct{}, 1),
		publicChan:  make(chan discovery.Event, 2),
		initialized: atomic.NewBool(false),
		option:      &option{},
	}

	return discovery
}

// ID returns the discovery provider id
func (d *Discovery) ID() string {
	return "kubernetes"
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

	// check the options
	if d.option.Provider == "" {
		d.option.Provider = d.ID()
	}

	return nil
}

// SetConfig registers the underlying discovery configuration
func (d *Discovery) SetConfig(meta discovery.Config) error {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()

	// first check whether the discovery provider is running
	if d.initialized.Load() {
		return discovery.ErrAlreadyInitialized
	}

	return d.setConfig(meta)
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

	// create the k8 config
	config, err := rest.InClusterConfig()
	// handle the error
	if err != nil {
		return errors.Wrap(err, "failed to get the in-cluster config of the kubernetes provider")
	}
	// create an instance of the k8 client set
	client, err := kubernetes.NewForConfig(config)
	// handle the error
	if err != nil {
		return errors.Wrap(err, "failed to create the kubernetes client api")
	}
	// set the k8 client
	d.client = client
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
	// set the initialized to false
	d.initialized = atomic.NewBool(false)
	// stop the watchers
	close(d.stopChan)
	// close the public channel
	close(d.publicChan)
	// return
	return nil
}

// Watch returns event based upon node lifecycle
func (d *Discovery) Watch(ctx context.Context) (<-chan discovery.Event, error) {
	// first check whether the actor system has started
	if !d.initialized.Load() {
		return nil, errors.New("kubernetes discovery engine not initialized")
	}
	// run the watcher
	go d.watchPods()
	return d.publicChan, nil
}

// DiscoverNodes returns a list of known nodes.
func (d *Discovery) DiscoverNodes() ([]*discovery.Node, error) {
	// first check whether the discovery provider is running
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}

	// let us create the pod labels map
	// TODO: make sure to document it on k8 discovery
	podLabels := map[string]string{
		"app.kubernetes.io/part-of":   d.option.ApplicationName,
		"app.kubernetes.io/component": d.option.ApplicationName, // TODO: redefine it
		"app.kubernetes.io/name":      d.option.ApplicationName,
	}

	// create a context
	ctx := context.Background()

	// List all the pods based on the filters we requested
	pods, err := d.client.CoreV1().Pods(d.option.NameSpace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(podLabels).String(),
	})
	// return an error when we cannot poll the pods
	if err != nil {
		return nil, err
	}

	nodes := make([]*discovery.Node, 0, pods.Size())
	// iterate the pods list and only the one that are running
MainLoop:
	for _, pod := range pods.Items {
		// create a variable copy of pod
		pod := pod
		// only consider running pods
		if pod.Status.Phase != corev1.PodRunning {
			continue MainLoop
		}
		// If there is a Ready condition available, we need that to be true.
		// If no ready condition is set, then we accept this pod regardless.
		for _, condition := range pod.Status.Conditions {
			// ignore pod that is not in ready state
			if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
				continue MainLoop
			}
		}
		// create a variable holding the node
		node := d.podToNode(&pod)
		// continue the loop when we did not find any node
		if node == nil {
			continue MainLoop
		}
		// add the node to the list of nodes
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// Close closes the provider
func (d *Discovery) Close() error {
	return nil
}

// setConfig sets the kubernetes discoConfig
func (d *Discovery) setConfig(config discovery.Config) (err error) {
	// create an instance of option
	option := new(option)
	// extract the namespace
	option.NameSpace, err = config.GetString(Namespace)
	// handle the error in case the namespace value is not properly set
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
	d.option = option
	return nil
}

// podToNode takes a kubernetes pod and returns a Node
func (d *Discovery) podToNode(pod *corev1.Pod) *discovery.Node {
	// If there is a Ready condition available, we need that to be true.
	// If no ready condition is set, then we accept this pod regardless.
	for _, condition := range pod.Status.Conditions {
		// ignore pod that is not in ready state
		if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
			return nil
		}
	}

	// create a variable holding the node
	var node *discovery.Node
	// define valid port names
	validPortNames := []string{GroupCachePortName}

	// iterate the pod containers and find the named port
	for i := 0; i < len(pod.Spec.Containers) && node == nil; i++ {
		// let us get the container
		container := pod.Spec.Containers[i]

		// iterate the container ports to set the join port
		for _, port := range container.Ports {
			// make sure we have the gossip and cluster port defined
			if !slices.Contains(validPortNames, port.Name) {
				// skip that port
				continue
			}

			// set the node
			node = discovery.NewNode(pod.GetName(), pod.Status.PodIP, uint32(port.ContainerPort), pod.Status.StartTime.Time.UnixMilli())
			break
		}
	}

	return node
}

// handlePodAdded is called when a new pod is added
func (d *Discovery) handlePodAdded(pod *corev1.Pod) {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()
	// ignore the pod when it is not running
	if pod.Status.Phase != corev1.PodRunning {
		return
	}
	// get node
	node := d.podToNode(pod)
	// continue the loop when we did not find any node
	if node == nil {
		return
	}
	// here we find a node let us raise the node registered event
	event := &discovery.NodeAdded{Node: node}
	// add to the channel
	d.publicChan <- event
}

// handlePodUpdated is called when a pod is updated
func (d *Discovery) handlePodUpdated(old *corev1.Pod, pod *corev1.Pod) {
	// ignore the pod when it is not running
	if pod.Status.Phase != corev1.PodRunning {
		return
	}
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()
	// grab the old node
	oldNode := d.podToNode(old)
	// get the new node
	node := d.podToNode(pod)
	// continue the loop when we did not find any node
	if node == nil {
		return
	}
	// here we find a node let us raise the node modified event
	event := &discovery.NodeModified{
		Node:    node,
		Current: oldNode,
	}
	// add to the channel
	d.publicChan <- event
}

// handlePodDeleted is called when pod is deleted
func (d *Discovery) handlePodDeleted(pod *corev1.Pod) {
	// acquire the lock
	d.mu.Lock()
	// release the lock
	defer d.mu.Unlock()
	// get the new node
	node := d.podToNode(pod)
	// here we find a node let us raise the node removed event
	event := &discovery.NodeRemoved{Node: node}
	// add to the channel
	d.publicChan <- event
}

// watchPods keeps a watch on kubernetes pods activities and emit
// respective event when needed
func (d *Discovery) watchPods() {
	// TODO: make sure to document it on k8 discovery
	podLabels := map[string]string{
		"app.kubernetes.io/part-of":   d.option.ApplicationName,
		"app.kubernetes.io/component": d.option.ApplicationName, // TODO: redefine it
		"app.kubernetes.io/name":      d.option.ApplicationName,
	}
	// create the k8 informer factory
	factory := informers.NewSharedInformerFactoryWithOptions(
		d.client,
		10*time.Minute, // TODO make it configurable
		informers.WithNamespace(d.option.NameSpace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.LabelSelector = labels.SelectorFromSet(podLabels).String()
		}))
	// create the pods informer instance
	informer := factory.Core().V1().Pods().Informer()
	synced := false
	mux := &sync.RWMutex{}
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			mux.RLock()
			defer mux.RUnlock()
			if !synced {
				return
			}

			// Handler logic
			pod := obj.(*corev1.Pod)
			// handle the newly added pod
			d.handlePodAdded(pod)
		},
		UpdateFunc: func(current, node any) {
			mux.RLock()
			defer mux.RUnlock()
			if !synced {
				return
			}

			// Handler logic
			old := current.(*corev1.Pod)
			pod := node.(*corev1.Pod)
			d.handlePodUpdated(old, pod)
		},
		DeleteFunc: func(obj any) {
			mux.RLock()
			defer mux.RUnlock()
			if !synced {
				return
			}

			// Handler logic
			pod := obj.(*corev1.Pod)
			// handle the newly added pod
			d.handlePodDeleted(pod)
		},
	})
	if err != nil {
		return
	}

	// run the informer
	go informer.Run(d.stopChan)

	// wait for caches to sync
	isSynced := cache.WaitForCacheSync(d.stopChan, informer.HasSynced)
	mux.Lock()
	synced = isSynced
	mux.Unlock()

	// caches failed to sync
	if !synced {
		// TODO add a fatal logger instead of panic
		panic("caches failed to sync")
	}
}
