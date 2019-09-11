package scan

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/log"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/db"
)

const (
	StatusInitializing = "initializing"
	StatusGathering    = "gathering"
	StatusAnalyzing    = "analyzing"
	StatusFinished     = "finished"
	ContentScan        = "Content Scan"
	PathScan           = "Path Scan"
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

type Session struct {
	sync.Mutex

	Options  Options     `json:"-"`
	Out      *log.Logger `json:"-"`
	Stats    *Stats
	Findings []*Finding
	Store    *db.MysqlHandler
}

func (s *Session) Start() {
	s.InitStats()
	s.InitLogger()
	s.InitThreads()
	s.InitDB()
	// s.InitRouter()
}

func (s *Session) Finish() {
	s.Stats.FinishedAt = time.Now()
	s.Stats.Status = StatusFinished
	//closing db connectoin
	err := s.Store.CloseConnection()
	if err != nil {
		fmt.Println("Unable to close db connection: ", err)
	}
}

func (s *Session) InitStats() {
	if s.Stats != nil {
		return
	}
	s.Stats = &Stats{
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

func (s *Session) InitLogger() {
	s.Out = &log.Logger{}
	s.Out.SetDebug(*s.Options.Debug)
	s.Out.SetSilent(*s.Options.Silent)
}

func (s *Session) InitDB() {
	//Opening DB Connection
	s.Store = db.GetInstance()
	err := s.Store.OpenConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		fmt.Println("Unable to open db connection: ", err)
		// os.Exit(1)
	}
}

func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}
	runtime.GOMAXPROCS(*s.Options.Threads + 2) // thread count + main + web server
}

// func (s *Session) InitRouter() {
//   bind := fmt.Sprintf("%s:%d", *s.Options.BindAddress, *s.Options.Port)
//   s.Router = NewRouter(s)
//   go func(sess *Session) {
//     if err := sess.Router.Run(bind); err != nil {
//       sess.Out.Fatal("Error when starting web server: %s\n", err)
//     }
//   }(s)
// }

func (s *Session) AddFinding(finding *Finding) {
	s.Lock()
	defer s.Unlock()
	s.Findings = append(s.Findings, finding)
}

func (s *Session) SaveToFile(location string) error {
	sessionJson, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(location, sessionJson, 0644)
	if err != nil {
		return err
	}
	return nil
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

func ValidateNewSession(session *Session) error {
	// var err error

	// if session.Options, err = ParseOptions(); err != nil {
	//   return err
	// }

	if *session.Options.Save != "" && FileExists(*session.Options.Save) {
		return errors.New(fmt.Sprintf("File: %s already exists.", *session.Options.Save))
	}

	if *session.Options.Load != "" {
		if !FileExists(*session.Options.Load) {
			return errors.New(fmt.Sprintf("Session file %s does not exist or is not readable.", *session.Options.Load))
		}
		data, err := ioutil.ReadFile(*session.Options.Load)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &session); err != nil {
			return errors.New(fmt.Sprintf("Session file %s is corrupt or not generated this version.", *session.Options.Load))
		}
	}
	return nil
}

func PrintSessionStats(sess *Session) {
	sess.Out.Info("\nFindings....: %d\n", sess.Stats.Findings)
	sess.Out.Info("Files.......: %d\n", sess.Stats.Files)
	sess.Out.Info("Commits.....: %d\n", sess.Stats.Commits)
	sess.Out.Info("Repositories: %d\n", sess.Stats.Repositories)
	sess.Out.Info("Targets.....: %d\n\n", sess.Stats.Targets)
}
