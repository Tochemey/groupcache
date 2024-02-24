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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		meta := Options{
			"key-1": "value-1",
			"key-2": "value-2",
		}
		key := "key-1"
		actual, err := meta.GetString(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
		expected := "value-1"
		assert.Equal(t, expected, actual)
	})
	t.Run("With key not found", func(t *testing.T) {
		meta := Options{
			"key-1": "value-1",
			"key-2": "value-2",
		}
		key := "key-3"
		actual, err := meta.GetString(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "key=key-3 not found")
		assert.Empty(t, actual)
	})
	t.Run("With key value not of a type string", func(t *testing.T) {
		meta := Options{
			"key-1": "value-1",
			"key-2": 13,
		}
		key := "key-2"
		actual, err := meta.GetString(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "the key value is not a string")
		assert.Empty(t, actual)
	})
}

func TestGetMapString(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		meta := Options{
			"key-1": map[string]string{
				"key-11": "value-11",
				"key-12": "value-12",
			},
		}
		key := "key-1"
		actual, err := meta.GetMapString(key)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		expected := map[string]string{
			"key-11": "value-11",
			"key-12": "value-12",
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("With key not found", func(t *testing.T) {
		meta := Options{
			"key-1": map[string]string{
				"key-11": "value-11",
				"key-12": "value-12",
			},
		}
		key := "key-3"
		actual, err := meta.GetMapString(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "key=key-3 not found")
		assert.Empty(t, actual)
	})
	t.Run("With key value not of a type map[string]string", func(t *testing.T) {
		meta := Options{
			"key-2": 13,
		}
		key := "key-2"
		actual, err := meta.GetMapString(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "the key value is not a map[string]string")
		assert.Empty(t, actual)
	})
}

func TestGetInt(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		meta := Options{
			"key-1": 20,
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetInt(key)
		assert.NoError(t, err)
		assert.NotZero(t, actual)
		expected := 20
		assert.EqualValues(t, expected, actual)
	})
	t.Run("With key not found", func(t *testing.T) {
		meta := Options{
			"key-1": 20,
			"key-2": 30,
		}
		key := "key-3"
		actual, err := meta.GetInt(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "key=key-3 not found")
		assert.Zero(t, actual)
	})
	t.Run("With key value not an int", func(t *testing.T) {
		meta := Options{
			"key-1": "a",
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetInt(key)
		assert.Error(t, err)
		assert.Zero(t, actual)
	})
	t.Run("With key value a string int", func(t *testing.T) {
		meta := Options{
			"key-1": "20",
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetInt(key)
		assert.NoError(t, err)
		assert.NotZero(t, actual)
		expected := 20
		assert.EqualValues(t, expected, actual)
	})
}

func TestGetBool(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		meta := Options{
			"key-1": true,
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetBool(key)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.True(t, *actual)
	})
	t.Run("With key not found", func(t *testing.T) {
		meta := Options{
			"key-1": 20,
			"key-2": 30,
		}
		key := "key-3"
		actual, err := meta.GetBool(key)
		assert.Error(t, err)
		assert.EqualError(t, err, "key=key-3 not found")
		assert.Nil(t, actual)
	})
	t.Run("With key value not an boolean", func(t *testing.T) {
		meta := Options{
			"key-1": "a",
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetBool(key)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
	t.Run("With key value a string boolean", func(t *testing.T) {
		meta := Options{
			"key-1": "TRUE",
			"key-2": 30,
		}
		key := "key-1"
		actual, err := meta.GetBool(key)
		assert.NoError(t, err)
		assert.True(t, *actual)
	})
}
