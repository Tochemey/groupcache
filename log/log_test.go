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

package log

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	t.Run("With Debug log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(DebugLevel, buffer)
		// assert Debug log
		logger.Debug("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, DebugLevel.String(), lvl)
		require.Equal(t, DebugLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Debugf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, DebugLevel.String(), lvl)
		require.Equal(t, DebugLevel, logger.LogLevel())
	})
	t.Run("With Info log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(InfoLevel, buffer)
		// assert Debug log
		logger.Debug("test debug")
		require.Empty(t, buffer.String())
	})
	t.Run("With Error log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(ErrorLevel, buffer)
		// assert Debug log
		logger.Debug("test debug")
		require.Empty(t, buffer.String())
	})
}

func TestInfo(t *testing.T) {
	t.Run("With Info log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(InfoLevel, buffer)
		// assert Debug log
		logger.Info("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, InfoLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Infof("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, InfoLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())
	})
	t.Run("With Debug log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(DebugLevel, buffer)
		// assert Debug log
		logger.Info("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, InfoLevel.String(), lvl)

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Infof("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, InfoLevel.String(), lvl)
	})
	t.Run("With Error log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(ErrorLevel, buffer)
		// assert Debug log
		logger.Info("test debug")
		require.Empty(t, buffer.String())
	})
}

func TestWarn(t *testing.T) {
	t.Run("With Warn log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(WarningLevel, buffer)
		// assert Debug log
		logger.Warn("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, WarningLevel.String(), lvl)
		require.Equal(t, WarningLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Warnf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, WarningLevel.String(), lvl)
		require.Equal(t, WarningLevel, logger.LogLevel())
	})
	t.Run("With Debug log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(DebugLevel, buffer)
		// assert Debug log
		logger.Warn("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, WarningLevel.String(), lvl)

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Warnf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, WarningLevel.String(), lvl)
	})
	t.Run("With Error log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(ErrorLevel, buffer)
		// assert Debug log
		logger.Warn("test debug")
		require.Empty(t, buffer.String())
	})
}

func TestError(t *testing.T) {
	t.Run("With the Error log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(InfoLevel, buffer)
		// assert Debug log
		logger.Error("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Errorf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())
	})
	t.Run("With the Debug log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(DebugLevel, buffer)
		// assert Debug log
		logger.Error("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, DebugLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Errorf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, DebugLevel, logger.LogLevel())
	})
	t.Run("With the Info log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(InfoLevel, buffer)
		// assert Debug log
		logger.Error("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Errorf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, InfoLevel, logger.LogLevel())
	})
	t.Run("With the Warn log level", func(t *testing.T) {
		// create a bytes buffer that implements an io.Writer
		buffer := new(bytes.Buffer)
		// create an instance of Log
		logger := New(WarningLevel, buffer)
		// assert Debug log
		logger.Error("test debug")
		expected := "test debug"
		actual, err := extractMessage(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		lvl, err := extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, WarningLevel, logger.LogLevel())

		// reset the buffer
		buffer.Reset()
		// assert Debug log
		name := "world"
		logger.Errorf("hello %s", name)
		actual, err = extractMessage(buffer.Bytes())
		require.NoError(t, err)
		expected = "hello world"
		require.Equal(t, expected, actual)

		lvl, err = extractLevel(buffer.Bytes())
		require.NoError(t, err)
		require.Equal(t, ErrorLevel.String(), lvl)
		require.Equal(t, WarningLevel, logger.LogLevel())
	})
}

func TestPanic(t *testing.T) {
	// create a bytes buffer that implements an io.Writer
	buffer := new(bytes.Buffer)
	// create an instance of Log
	logger := New(PanicLevel, buffer)
	// assert Debug log
	assert.Panics(t, func() {
		logger.Panic("test debug")
	})
}

func extractMessage(bytes []byte) (string, error) {
	// a map container to decode the JSON structure into
	c := make(map[string]json.RawMessage)

	// unmarshal JSON
	if err := json.Unmarshal(bytes, &c); err != nil {
		return "", err
	}
	for k, v := range c {
		if k == "msg" {
			return strconv.Unquote(string(v))
		}
	}

	return "", nil
}

func extractLevel(bytes []byte) (string, error) {
	// a map container to decode the JSON structure into
	c := make(map[string]json.RawMessage)

	// unmarshal JSON
	if err := json.Unmarshal(bytes, &c); err != nil {
		return "", err
	}
	for k, v := range c {
		if k == "level" {
			return strconv.Unquote(string(v))
		}
	}

	return "", nil
}
