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
	"fmt"
	"net"
	"strconv"
)

// Node defines the discovered node
type Node struct {
	// name specifies the discovered node's name
	name string
	// host specifies the discovered node's Host
	host string
	// port
	port int32
}

// Name returns the Node name
func (x Node) Name() string {
	return x.name
}

// Host returns the peer host
func (x Node) Host() string {
	return x.host
}

// Port returns the Node port number
func (x Node) Port() int32 {
	return x.port
}

// NewNode creates an instance of Node
func NewNode(name, host string, port int32) *Node {
	return &Node{
		name: name,
		host: host,
		port: port,
	}
}

// Address returns the Node address
func (x Node) Address() string {
	return net.JoinHostPort(x.host, strconv.Itoa(int(x.port)))
}

// PeerURL returns the PeerURL
func (x Node) PeerURL() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(x.host, strconv.Itoa(int(x.port))))
}
