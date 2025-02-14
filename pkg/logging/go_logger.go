/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains a logger that uses the Go `log` package.

package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"io"
)

// GoLoggerBuilder contains the configuration and logic needed to build a logger that uses the Go
// `log` package. Don't create instances of this type directly, use the NewGoLoggerBuilder function
// instead.
type GoLoggerBuilder struct {
	debugEnabled bool
	infoEnabled  bool
	warnEnabled  bool
	errorEnabled bool
	logFile string
}

// GoLogger is a logger that uses the Go `log` package.
type GoLogger struct {
	debugEnabled bool
	infoEnabled  bool
	warnEnabled  bool
	errorEnabled bool
	logFileEnabled bool
	logFile string
	logFileToClose *os.File
}

// NewGoLoggerBuilder creates a builder that knows how to build a logger that uses the Go `log`
// package. By default these loggers will have enabled the information, warning and error levels
func NewGoLoggerBuilder() *GoLoggerBuilder {
	// Allocate the object:
	builder := new(GoLoggerBuilder)

	// Set default values:
	builder.debugEnabled = false
	builder.infoEnabled = true
	builder.warnEnabled = true
	builder.errorEnabled = true
        builder.logFile = ""

	return builder
}

// Debug enables or disables the debug level.
func (b *GoLoggerBuilder) Debug(flag bool) *GoLoggerBuilder {
	b.debugEnabled = flag
	return b
}

// Info enables or disables the information level.
func (b *GoLoggerBuilder) Info(flag bool) *GoLoggerBuilder {
	b.infoEnabled = flag
	return b
}

// Warn enables or disables the warning level.
func (b *GoLoggerBuilder) Warn(flag bool) *GoLoggerBuilder {
	b.warnEnabled = flag
	return b
}

// Error enables or disables the error level.
func (b *GoLoggerBuilder) Error(flag bool) *GoLoggerBuilder {
	b.errorEnabled = flag
	return b
}

// Set log file location
func (b *GoLoggerBuilder) LogFile(flag string) *GoLoggerBuilder {
        b.logFile = flag
        return b
}

// Build creates a new logger using the configuration stored in the builder.
func (b *GoLoggerBuilder) Build() (logger *GoLogger, err error) {
	// Allocate and populate the object:
	logger = new(GoLogger)
	logger.debugEnabled = b.debugEnabled
	logger.infoEnabled = b.infoEnabled
	logger.warnEnabled = b.warnEnabled
	logger.errorEnabled = b.errorEnabled
	logger.logFile = b.logFile

	if b.logFile != "" {
                lFile, err := os.OpenFile(b.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
                if err != nil {
                        log.Fatalf("Error opening log-file for writing: %v\n", err)
                }
		logger.logFileToClose = lFile
                logger.SetOutput(lFile)
        }

	return
}

// DebugEnabled returns true iff the debug level is enabled.
func (l *GoLogger) DebugEnabled() bool {
	return l.debugEnabled
}

// InfoEnabled returns true iff the information level is enabled.
func (l *GoLogger) InfoEnabled() bool {
	return l.infoEnabled
}

// WarnEnabled returns true iff the warning level is enabled.
func (l *GoLogger) WarnEnabled() bool {
	return l.warnEnabled
}

// ErrorEnabled returns true iff the error level is enabled.
func (l *GoLogger) ErrorEnabled() bool {
	return l.errorEnabled
}

// LogFileEnabled returns true iff there is a log file set
func (l *GoLogger) LogFileEnabled() bool {
        if l.logFile != "" {
                return true
        }
        return false
}

// Debug sends to the log a debug message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GoLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	if l.debugEnabled {
		format = appendHeader(Debug, format)
		msg := fmt.Sprintf(format, args...)
		// #nosec G104
		log.Output(1, msg)
	}
}

// Info sends to the log an information message formatted using the fmt.Sprintf function and the
// given format and arguments.
func (l *GoLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.infoEnabled {
		format = appendHeader(Info, format)
		msg := fmt.Sprintf(format, args...)
		// #nosec G104
		log.Output(1, msg)
	}
}

// Warn sends to the log a warning message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GoLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.warnEnabled {
		format = appendHeader(Warning, format)
		msg := fmt.Sprintf(format, args...)
		// #nosec G104
		log.Output(1, msg)
	}
}

// Error sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GoLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.errorEnabled {
		format = appendHeader(Error, format)
		msg := fmt.Sprintf(format, args...)
		// #nosec G104
		if l.LogFileEnabled() {
                        mw := io.MultiWriter(log.Writer(), os.Stderr)
                        log.SetOutput(mw)
                }
                log.Output(1, msg)
                if l.LogFileEnabled() {
                        log.SetOutput(l.logFileToClose)
                }
	}
}

// Fatal sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments. After that it will os.Exit(1)
// This level is always enabled
func (l *GoLogger) Fatal(ctx context.Context, format string, args ...interface{}) {
	format = appendHeader(Fatal, format)
	msg := fmt.Sprintf(format, args...)
	// #nosec G104
	if l.LogFileEnabled() {
                mw := io.MultiWriter(log.Writer(), os.Stderr)
                log.SetOutput(mw)
        }
	log.Output(1, msg)
	l.DeferClose()
	os.Exit(1)
}

// Set the output of the logger to a log file
func (l *GoLogger) SetOutput(w io.Writer) {
        log.SetOutput(w)
}

// If we are writing to a log file, close it
func (l *GoLogger) DeferClose() {
        if l.LogFileEnabled() {
                l.logFileToClose.Close()
        }
}
