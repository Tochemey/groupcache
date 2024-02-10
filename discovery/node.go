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
	port uint32
	// birthday
	birthday int64
}

// Name returns the Node name
func (x Node) Name() string {
	return x.name
}

// Self returns the peer host
func (x Node) Host() string {
	return x.host
}

// Port returns the Node port number
func (x Node) Port() uint32 {
	return x.port
}

// Birthday returns the Node date of birth
func (x Node) Birthday() int64 {
	return x.birthday
}

// NewNode creates an instance of Node
func NewNode(name, host string, port uint32, dob int64) *Node {
	return &Node{
		name:     name,
		host:     host,
		port:     port,
		birthday: dob,
	}
}

// Address returns the Node address
func (x Node) Address() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(x.host, strconv.Itoa(int(x.port))))
}
