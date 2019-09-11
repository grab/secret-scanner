package stats

import (
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/log"
	"sync"
	"time"
)

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

func (s *Stats) IncrementTargets() {
	s.Lock()
	defer s.Unlock()
	s.Targets++
}

func (s *Stats) IncrementRepositories() {
	s.Lock()
	defer s.Unlock()
	s.Repositories++
}

func (s *Stats) IncrementCommits() {
	s.Lock()
	defer s.Unlock()
	s.Commits++
}

func (s *Stats) IncrementFiles() {
	s.Lock()
	defer s.Unlock()
	s.Files++
}

func (s *Stats) IncrementFindings() {
	s.Lock()
	defer s.Unlock()
	s.Findings++
}

func (s *Stats) UpdateProgress(current int, total int) {
	s.Lock()
	defer s.Unlock()
	if current >= total {
		s.Progress = 100.0
	} else {
		s.Progress = (float64(current) * float64(100)) / float64(total)
	}
}

func (s *Stats) PrintStats(logger *log.Logger) {
	logger.Info("\nFindings....: %d\n", s.Findings)
	logger.Info("Files.......: %d\n", s.Files)
	logger.Info("Commits.....: %d\n", s.Commits)
	logger.Info("Repositories: %d\n", s.Repositories)
	logger.Info("Targets.....: %d\n\n", s.Targets)
}
