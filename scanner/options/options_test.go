/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package options

import (
	"testing"
)

func TestOptions_ParseScanTargets(t *testing.T) {
	sampleTargets := []string{"123", "456", "7899"}
	scanTargetStr := "123, 456 , 7899"
	options := Options{ScanTarget: &scanTargetStr}
	targets := options.ParseScanTargets()
	if numT := len(targets); numT != 3 {
		t.Errorf("Want 3, got %v", numT)
		return
	}
	for i, target := range targets {
		if target != sampleTargets[i] {
			t.Errorf(`Want "%s", got %s`, sampleTargets[i], target)
		}
	}
}
