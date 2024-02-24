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

package discovery

// Discovery defines the cluster service discovery
type Discovery struct {
	// provider specifies the discovery provider
	provider Provider
	// options specifies the discovery options
	options Options
}

// New creates an instance of Discovery
func New(provider Provider, options Options) *Discovery {
	return &Discovery{
		provider: provider,
		options:  options,
	}
}

// Provider returns the service discovery provider
func (s Discovery) Provider() Provider {
	return s.provider
}

// Options returns the service discovery options
func (s Discovery) Options() Options {
	return s.options
}
