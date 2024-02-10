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
	"fmt"
	"io"
	golog "log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultLogger represents the default Log to use
// This Log wraps zerolog under the hood
var DefaultLogger = New(DebugLevel, os.Stdout)
var DiscardLogger = New(DebugLevel, io.Discard)

// Log implements Logger interface with the underlying zap as
// the underlying logging library
type Log struct {
	logger  *zap.Logger
	outputs []io.Writer
}

// enforce compilation and linter error
var _ Logger = &Log{}

// New creates an instance of Log
// New creates an instance of Log
func New(level Level, writers ...io.Writer) *Log {
	// create the zap Log configuration
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	// create the zap log core
	var core zapcore.Core

	// create the list of writers
	syncWriters := make([]zapcore.WriteSyncer, len(writers))
	for i, writer := range writers {
		syncWriters[i] = zapcore.AddSync(writer)
	}

	// set the log level
	switch level {
	case InfoLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.InfoLevel,
		)
	case DebugLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.DebugLevel,
		)
	case WarningLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.WarnLevel,
		)
	case ErrorLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.ErrorLevel,
		)
	case PanicLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.PanicLevel,
		)
	case FatalLevel:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.FatalLevel,
		)
	default:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zap.CombineWriteSyncers(syncWriters...),
			zapcore.DebugLevel,
		)
	}
	// get the zap Log
	zapLogger := zap.New(core)
	// set the global logger
	zap.ReplaceGlobals(zapLogger)
	// create the instance of Log and returns it
	return &Log{
		logger:  zapLogger,
		outputs: writers,
	}
}

// Debug starts a message with debug level
func (l *Log) Debug(v ...any) {
	defer l.logger.Sync()
	l.logger.Debug(fmt.Sprint(v...))
}

// Debugf starts a message with debug level
func (l *Log) Debugf(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Debug(fmt.Sprintf(format, v...))
}

// Panic starts a new message with panic level. The panic() function
// is called which stops the ordinary flow of a goroutine.
func (l *Log) Panic(v ...any) {
	defer l.logger.Sync()
	l.logger.Panic(fmt.Sprint(v...))
}

// Panicf starts a new message with panic level. The panic() function
// is called which stops the ordinary flow of a goroutine.
func (l *Log) Panicf(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Panic(fmt.Sprintf(format, v...))
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called which terminates the program immediately.
func (l *Log) Fatal(v ...any) {
	defer l.logger.Sync()
	l.logger.Fatal(fmt.Sprint(v...))
}

// Fatalf starts a new message with fatal level. The os.Exit(1) function
// is called which terminates the program immediately.
func (l *Log) Fatalf(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

// Error starts a new message with error level.
func (l *Log) Error(v ...any) {
	defer l.logger.Sync()
	l.logger.Error(fmt.Sprint(v...))
}

// Errorf starts a new message with error level.
func (l *Log) Errorf(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Error(fmt.Sprintf(format, v...))
}

// Warn starts a new message with warn level
func (l *Log) Warn(v ...any) {
	defer l.logger.Sync()
	l.logger.Warn(fmt.Sprint(v...))
}

// Warnf starts a new message with warn level
func (l *Log) Warnf(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Warn(fmt.Sprintf(format, v...))
}

// Info starts a message with info level
func (l *Log) Info(v ...any) {
	defer l.logger.Sync()
	l.logger.Info(fmt.Sprint(v...))
}

// Infof starts a message with info level
func (l *Log) Infof(format string, v ...any) {
	defer l.logger.Sync()
	l.logger.Info(fmt.Sprintf(format, v...))
}

// LogLevel returns the log level that is used
func (l *Log) LogLevel() Level {
	switch l.logger.Level() {
	case zapcore.FatalLevel:
		return FatalLevel
	case zapcore.PanicLevel:
		return PanicLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.WarnLevel:
		return WarningLevel
	default:
		return InvalidLevel
	}
}

// LogOutput returns the log output that is set
func (l *Log) LogOutput() []io.Writer {
	return l.outputs
}

// StdLogger returns the standard logger associated to the logger
func (l *Log) StdLogger() *golog.Logger {
	stdlogger, _ := zap.NewStdLogAt(l.logger, l.logger.Level())
	redirect, _ := zap.RedirectStdLogAt(l.logger, l.logger.Level())
	defer redirect()
	return stdlogger
}

// With set field and key in the logger
func (l *Log) With(key string, value any) Logger {
	logger := l.logger.With(zap.Any(key, value))
	l.logger = logger
	return l
}
