/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

const (
	// FATAL ...
	FATAL = 5
	// ERROR ...
	ERROR = 4
	// WARN ...
	WARN = 3
	// IMPORTANT ...
	IMPORTANT = 2
	// INFO ...
	INFO = 1
	// DEBUG ...
	DEBUG = 0
)

// LogColors defines message output colors
var LogColors = map[int]*color.Color{
	FATAL:     color.New(color.FgRed).Add(color.Bold),
	ERROR:     color.New(color.FgRed),
	WARN:      color.New(color.FgYellow),
	IMPORTANT: color.New(color.Bold),
	DEBUG:     color.New(color.FgCyan).Add(color.Faint),
}

// Logger ...
type Logger struct {
	sync.Mutex

	debug  bool
	silent bool
}

// SetSilent sets logger silent option
func (l *Logger) SetSilent(s bool) {
	l.silent = s
}

// SetDebug sets logger debug option
func (l *Logger) SetDebug(d bool) {
	l.debug = d
}

// Log logs a message
func (l *Logger) Log(level int, format string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if level == DEBUG && l.debug == false {
		return
	} else if level < ERROR && l.silent == true {
		return
	}

	if c, ok := LogColors[level]; ok {
		c.Printf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

// Fatal ...
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}

// Error ...
func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

// Warn ...
func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(WARN, format, args...)
}

// Important ...
func (l *Logger) Important(format string, args ...interface{}) {
	l.Log(IMPORTANT, format, args...)
}

// Info ...
func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

// Debug ...
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}
