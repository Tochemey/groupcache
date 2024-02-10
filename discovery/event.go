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

// Event contract
type Event interface {
	IsEvent()
}

// NodeAdded discovery lifecycle event
type NodeAdded struct {
	// Node specifies the added node
	Node *Node
}

func (e NodeAdded) IsEvent() {}

// NodeRemoved discovery lifecycle event
type NodeRemoved struct {
	// Node specifies the removed node
	Node *Node
}

func (e NodeRemoved) IsEvent() {}

// NodeModified discovery lifecycle event
type NodeModified struct {
	// Node specifies the modified node
	Node *Node
	// Current specifies the existing nde
	Current *Node
}

func (e NodeModified) IsEvent() {}
