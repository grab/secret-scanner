/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package filehandler

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFileExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "ss-test-")
	if err != nil {
		t.Errorf("Cannot create temp. dir.: %v", err)
		return
	}

	f, err := os.Create(path.Join(tempDir, "ss-test-file"))
	if err != nil {
		t.Errorf("Cannot create file: %v", err)
		return
	}

	if !FileExists(f.Name()) {
		t.Errorf("Want file %s exists, got not exists", f.Name())
		return
	}

	cleanup(tempDir)

	if FileExists(f.Name()) {
		t.Errorf("Want file %s not exists, got exists", f.Name())
	}
}

func cleanup(path string) {
	_ = os.RemoveAll(path)
}
