/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package stats

import (
	"sync"
	"time"

	"github.com/grab/secret-scanner/common/log"
)

// Stats holds info about the scan status
type Stats struct {
	sync.Mutex

	StartedAt    time.Time
	FinishedAt   time.Time
	Status       string
	Progress     float64
	Targets      int
	Repositories int
	Commits      int
	Files        int
	Findings     int
}

// IncrementTargets increase the target count by 1
func (s *Stats) IncrementTargets() {
	s.Lock()
	defer s.Unlock()
	s.Targets++
}

// IncrementRepositories increase the repo count by 1
func (s *Stats) IncrementRepositories() {
	s.Lock()
	defer s.Unlock()
	s.Repositories++
}

// IncrementCommits increase commit count by 1
func (s *Stats) IncrementCommits() {
	s.Lock()
	defer s.Unlock()
	s.Commits++
}

// IncrementFiles increase file count by 1
func (s *Stats) IncrementFiles() {
	s.Lock()
	defer s.Unlock()
	s.Files++
}

// IncrementFindings increase finding count by 1
func (s *Stats) IncrementFindings() {
	s.Lock()
	defer s.Unlock()
	s.Findings++
}

// UpdateProgress updates the progress percentage
func (s *Stats) UpdateProgress(current int, total int) {
	s.Lock()
	defer s.Unlock()
	if current >= total {
		s.Progress = 100.0
	} else {
		s.Progress = (float64(current) * float64(100)) / float64(total)
	}
}

// PrintStats prints the stat info
func (s *Stats) PrintStats(logger *log.Logger) {
	logger.Info("\nFindings....: %d\n", s.Findings)
	logger.Info("Files.......: %d\n", s.Files)
	logger.Info("Commits.....: %d\n", s.Commits)
	logger.Info("Repositories: %d\n", s.Repositories)
	logger.Info("Targets.....: %d\n\n", s.Targets)
}
