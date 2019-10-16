/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package log

import (
	"sync"
	"testing"
)

func TestLogger_SetSilent(t *testing.T) {
	logger := &Logger{
		Mutex:  sync.Mutex{},
		debug:  false,
		silent: false,
	}

	logger.SetSilent(true)

	if !logger.silent {
		t.Errorf("Want true, got false")
	}
}

func TestLogger_SetDebug(t *testing.T) {
	logger := &Logger{
		Mutex:  sync.Mutex{},
		debug:  false,
		silent: false,
	}

	logger.SetDebug(true)

	if !logger.debug {
		t.Errorf("Want true, got false")
	}
}
