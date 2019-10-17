/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package signatures

const (
	// TypeSimple ...
	TypeSimple = "simple"

	// TypePattern ...
	TypePattern = "pattern"

	// PartExtension ...
	PartExtension = "extension"

	// PartFilename ...
	PartFilename = "filename"

	// PartPath ...
	PartPath = "path"

	// PartContent ...
	PartContent = "content"
)

var skippableExtensions = []string{".exe", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf", ".zip", ".tar.gz", ".ttf", ".lock"}
var skippablePathIndicators = []string{"node_modules/", "vendor/", "bin/"}
var skippableTestPaths = []string{"test", "_spec", "fixture", "mock", "stub", "fake", "demo", "sample"}
