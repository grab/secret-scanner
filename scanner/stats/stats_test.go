/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package stats

import (
	"testing"
	"time"
)

func TestStats_IncrementTargets(t *testing.T) {
	st := createNewStat()

	st.IncrementTargets()
	if st.Targets != 1 {
		t.Errorf("Want 1, got %v", st.Targets)
	}
}

func TestStats_IncrementRepositories(t *testing.T) {
	st := createNewStat()

	st.IncrementRepositories()
	if st.Repositories != 1 {
		t.Errorf("Want 1, got %v", st.Repositories)
	}
}

func TestStats_IncrementFiles(t *testing.T) {
	st := createNewStat()

	st.IncrementFiles()
	if st.Files != 1 {
		t.Errorf("Want 1, got %v", st.Files)
	}
}

func TestStats_IncrementFindings(t *testing.T) {
	st := createNewStat()

	st.IncrementFindings()
	if st.Findings != 1 {
		t.Errorf("Want 1, got %v", st.Findings)
	}
}

func TestStats_UpdateProgress(t *testing.T) {
	st := createNewStat()

	st.UpdateProgress(0, 10)
	if st.Progress != float64(0) {
		t.Errorf("Want 0, got %v", st.Progress)
	}

	st.UpdateProgress(5, 10)
	if st.Progress != float64(50) {
		t.Errorf("Want 50, got %v", st.Progress)
	}

	st.UpdateProgress(11, 10)
	if st.Progress != float64(100) {
		t.Errorf("Want 100, got %v", st.Progress)
	}
}

func createNewStat() *Stats {
	return &Stats{
		StartedAt:    time.Now(),
		Progress:     0.0,
		Targets:      0,
		Repositories: 0,
		Commits:      0,
		Files:        0,
		Findings:     0,
	}
}
