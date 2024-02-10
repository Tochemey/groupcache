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

import (
	"os"

	"github.com/caarlos0/env/v10"
)

type selfConfig struct {
	GroupPort uint32 `env:"GROUP_PORT"`
	Name      string `env:"NODE_NAME" envDefault:""`
	Host      string `env:"NODE_IP" envDefault:""`
}

// Self returns the Node where the discovery provider is running
func Self() (*Node, error) {
	// load the host node configuration
	cfg := &selfConfig{}
	opts := env.Options{RequiredIfNoDef: true, UseFieldNameByDefault: false}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}
	// check for empty host and name
	if cfg.Host == "" {
		// let us perform a host lookup
		host, _ := os.Hostname()
		// set the host
		cfg.Host = host
	}

	// set the name as host if it is empty
	if cfg.Name == "" {
		cfg.Name = cfg.Host
	}

	// create the host node
	return &Node{
		name: cfg.Name,
		host: cfg.Host,
		port: cfg.GroupPort,
	}, nil
}
