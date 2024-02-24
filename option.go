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

package groupcache

import (
	"time"

	"github.com/tochemey/groupcache/v2/consistenthash"
	"github.com/tochemey/groupcache/v2/log"
)

type Option interface {
	// Apply sets the Option value of a config.
	Apply(node *Node)
}

var _ Option = OptionFunc(nil)

// OptionFunc implements the Option interface.
type OptionFunc func(node *Node)

// Apply applies the Node's option
func (f OptionFunc) Apply(node *Node) {
	f(node)
}

// WithReplicaCount sets the total number of replicas count
func WithReplicaCount(count int) Option {
	return OptionFunc(func(cl *Node) {
		cl.replicaCount = count
	})
}

// WithLogger sets the logger
func WithLogger(logger log.Logger) Option {
	return OptionFunc(func(node *Node) {
		node.logger = logger
	})
}

// WithHashFn sets the custom hash function
func WithHashFn(hashFn consistenthash.Hash) Option {
	return OptionFunc(func(node *Node) {
		node.hashFn = hashFn
	})
}

// WithShutdownTimeout sets the shutdown timeout
func WithShutdownTimeout(timeout time.Duration) Option {
	return OptionFunc(func(node *Node) {
		node.shutdownTimeout = timeout
	})
}

// WithBasePath set the base path
func WithBasePath(path string) Option {
	return OptionFunc(func(node *Node) {
		node.basePath = path
	})
}
