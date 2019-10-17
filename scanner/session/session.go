/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/history"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/findings"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/filehandler"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/log"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/options"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/signatures"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/stats"
)

// Session contains fields describing a scan session
type Session struct {
	sync.Mutex

	Options      options.Options `json:"-"`
	Out          *log.Logger     `json:"-"`
	Stats        *stats.Stats
	Findings     []*findings.Finding
	Repositories []*gitprovider.Repository
	Signatures   []signatures.Signature `json:"-"`
	HistoryStore *history.JSONFileStore
}

// Initialize inits a scan session
func (s *Session) Initialize(options options.Options) {
	s.Options = options
	s.InitHistoryStoreOrFail(*options.HistoryStoreFilePath)
	s.InitLogger()
	s.InitStats()
	s.InitThreads()
	s.Signatures = signatures.LoadSignatures()
}

// End end a scan session
func (s *Session) End() {
	s.Stats.FinishedAt = time.Now()
	s.Stats.Status = StatusFinished
	s.HistoryStore.Close()
}

// InitHistoryStoreOrFail inits a history storage
func (s *Session) InitHistoryStoreOrFail(filepath string) {
	s.HistoryStore = &history.JSONFileStore{}

	if filepath == "" {
		defaultPath, err := s.HistoryStore.GetDefaultStorePath()
		if err != nil {
			fmt.Println(fmt.Sprintf("Unable to create default history file path: %v", err))
			os.Exit(1)
		}
		filepath = defaultPath
	}

	err := s.HistoryStore.Initialize(filepath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to initialize HistoryStore: %v", err))
		os.Exit(1)
	}
}

// InitLogger inits a logger
func (s *Session) InitLogger() {
	s.Out = &log.Logger{}
	s.Out.SetDebug(*s.Options.Debug)
	s.Out.SetSilent(*s.Options.Silent)
}

// InitStats inits stats
func (s *Session) InitStats() {
	s.Stats = &stats.Stats{
		StartedAt:    time.Now(),
		Status:       StatusInitializing,
		Progress:     0.0,
		Targets:      0,
		Repositories: 0,
		Commits:      0,
		Files:        0,
		Findings:     0,
	}
}

// InitThreads inits threads
func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}
	runtime.GOMAXPROCS(*s.Options.Threads + 2) // thread count + main + web server
}

// AddFinding adds a finding
func (s *Session) AddFinding(finding *findings.Finding) {
	s.Lock()
	defer s.Unlock()
	s.Findings = append(s.Findings, finding)
}

// SaveToFile exports scan results to file
func (s *Session) SaveToFile(location string) error {
	// get absolute path
	absPath, err := filepath.Abs(location)
	if err != nil {
		return err
	}

	// session to json bytes
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// if exists write to file
	if filehandler.FileExists(absPath) {
		err = ioutil.WriteFile(absPath, sessionJSON, 0644)
		if err != nil {
			return err
		}
		return nil
	}

	// create dirs
	dirPath := path.Dir(absPath)
	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return err
	}

	// create file
	_, err = os.Create(absPath)
	if err != nil {
		return err
	}

	// write to file
	err = ioutil.WriteFile(absPath, sessionJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AddRepository adds a repo
func (s *Session) AddRepository(repository *gitprovider.Repository) {
	s.Lock()
	defer s.Unlock()
	for _, r := range s.Repositories {
		if repository.ID == r.ID {
			return
		}
	}
	s.Repositories = append(s.Repositories, repository)
}

// ValidateNewSession validates new session
func ValidateNewSession(session *Session) error {
	// var err error

	// if session.Options, err = ParseOptions(); err != nil {
	//   return err
	// }

	if *session.Options.Save != "" && filehandler.FileExists(*session.Options.Save) {
		return fmt.Errorf("file: %s already exists", *session.Options.Save)
	}

	if *session.Options.Load != "" {
		if !filehandler.FileExists(*session.Options.Load) {
			return fmt.Errorf("session file %s does not exist or is not readable", *session.Options.Load)
		}
		data, err := ioutil.ReadFile(*session.Options.Load)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &session); err != nil {
			return fmt.Errorf("session file %s is corrupt or not generated this version", *session.Options.Load)
		}
	}
	return nil
}
