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

// ServiceDiscovery defines the cluster service discovery
type ServiceDiscovery struct {
	// provider specifies the discovery provider
	provider Provider
	// config specifies the discovery config
	config Config
}

// NewServiceDiscovery creates an instance of ServiceDiscovery
func NewServiceDiscovery(provider Provider, config Config) *ServiceDiscovery {
	return &ServiceDiscovery{
		provider: provider,
		config:   config,
	}
}

// Provider returns the service discovery provider
func (s ServiceDiscovery) Provider() Provider {
	return s.provider
}

// Config returns the service discovery config
func (s ServiceDiscovery) Config() Config {
	return s.config
}
