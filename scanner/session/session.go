/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/grab/secret-scanner/scanner/state"

	"github.com/grab/secret-scanner/scanner/findings"

	"github.com/grab/secret-scanner/common/filehandler"
	"github.com/grab/secret-scanner/common/log"
	"github.com/grab/secret-scanner/scanner/gitprovider"
	"github.com/grab/secret-scanner/scanner/options"
	"github.com/grab/secret-scanner/scanner/signatures"
	"github.com/grab/secret-scanner/scanner/stats"
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
	StateStore   *state.JSONFileStore
}

// Initialize inits a scan session
func (s *Session) Initialize(options options.Options) {
	s.Options = options
	s.InitStateStoreOrFail("")
	s.InitLogger()
	s.InitStats()
	s.InitThreads()
	s.Signatures = signatures.LoadSignatures()
}

// End end a scan session
func (s *Session) End() {
	s.Stats.FinishedAt = time.Now()
	s.Stats.Status = StatusFinished
	s.StateStore.Close()
}

// InitStateStoreOrFail inits a history storage
func (s *Session) InitStateStoreOrFail(filepath string) {
	s.StateStore = &state.JSONFileStore{}

	if filepath == "" {
		defaultPath, err := s.StateStore.GetDefaultStorePath()
		if err != nil {
			fmt.Println(fmt.Sprintf("Unable to create default history file path: %v", err))
			os.Exit(1)
		}
		filepath = defaultPath
	}

	err := s.StateStore.Initialize(filepath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to initialize StateStore: %v", err))
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
func (s *Session) SaveToFile(location string) (string, error) {
	// get absolute path
	absPath, err := filepath.Abs(location)
	if err != nil {
		return "", err
	}

	// session to json bytes
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	// prettify JSON
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, sessionJSON, "", "\t")
	if err != nil {
		return "", err
	}

	// if exists write to file
	if filehandler.FileExists(absPath) {
		err = ioutil.WriteFile(absPath, prettyJSON.Bytes(), 0644)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	// create dirs
	dirPath := path.Dir(absPath)
	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return "", err
	}

	// create file
	_, err = os.Create(absPath)
	if err != nil {
		return "", err
	}

	// write to file
	err = ioutil.WriteFile(absPath, prettyJSON.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return absPath, nil
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

	if *session.Options.Report != "" && filehandler.FileExists(*session.Options.Report) {
		return fmt.Errorf("file: %s already exists", *session.Options.Report)
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
