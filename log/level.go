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

// Level specifies the log level
type Level int

const (
	// InfoLevel indicates Info log level.
	InfoLevel Level = iota
	// WarningLevel indicates Warning log level.
	WarningLevel
	// ErrorLevel indicates Error log level.
	ErrorLevel
	// FatalLevel indicates Fatal log level.
	FatalLevel
	// PanicLevel indicates Panic log level
	PanicLevel
	// DebugLevel indicates Debug log level
	DebugLevel
	InvalidLevel
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case WarningLevel:
		return "warn"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case ErrorLevel:
		return "error"
	default:
		return "" // FIXME: Surely we should blast here
	}
}
