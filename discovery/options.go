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
	"errors"
	"fmt"
	"strconv"
)

// Options represents the meta information to pass to the discovery engine
type Options map[string]any

// NewOptions initializes meta
func NewOptions() Options {
	return Options{}
}

// GetString returns the string value of a given key which value is a string
// If the key value is not a string then an error is return
func (m Options) GetString(key string) (string, error) {
	// let us check whether the given key is in the map
	val, ok := m[key]
	if !ok {
		return "", fmt.Errorf("key=%s not found", key)
	}
	// let us check the type of val
	switch x := val.(type) {
	case string:
		return x, nil
	default:
		return "", errors.New("the key value is not a string")
	}
}

// GetInt returns the int value of a given key which value is an integer
// If the key value is not an integer then an error is return
func (m Options) GetInt(key string) (int, error) {
	// let us check whether the given key is in the map
	val, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key=%s not found", key)
	}
	// let us check the type of val
	switch x := val.(type) {
	case int:
		return x, nil
	default:
		// maybe it is string integer
		return strconv.Atoi(val.(string))
	}
}

// GetBool returns the int value of a given key which value is a boolean
// If the key value is not a boolean then an error is return
func (m Options) GetBool(key string) (*bool, error) {
	// let us check whether the given key is in the map
	val, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("key=%s not found", key)
	}
	// let us check the type of val
	switch x := val.(type) {
	case bool:
		return &x, nil
	default:
		// parse the string value
		res, err := strconv.ParseBool(val.(string))
		// return the possible error
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
}

// GetMapString returns the map of string value of a given key which value is a map of string
// Map of string means that the map key value pair are both string
func (m Options) GetMapString(key string) (map[string]string, error) {
	// let us check whether the given key is in the map
	val, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("key=%s not found", key)
	}
	// assert the type of val
	switch x := val.(type) {
	case map[string]string:
		return x, nil
	default:
		return nil, errors.New("the key value is not a map[string]string")
	}
}
