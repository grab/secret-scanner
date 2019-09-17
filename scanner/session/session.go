package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/filehandler"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/log"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/options"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/signatures"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/stats"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/db"
)

type Session struct {
	sync.Mutex

	Options      options.Options `json:"-"`
	Out          *log.Logger     `json:"-"`
	Stats        *stats.Stats
	Findings     []*signatures.Finding
	Store        *db.MysqlHandler
	Repositories []*gitprovider.Repository
}

func (s *Session) Initialize(options options.Options) {
	s.Options = options
	s.InitLogger()
	s.InitDB()
	s.InitStats()
	s.InitThreads()
}

func (s *Session) End() {
	s.Stats.FinishedAt = time.Now()
	s.Stats.Status = StatusFinished
	//closing db connectoin
	err := s.Store.CloseConnection()
	if err != nil {
		fmt.Println("Unable to close db connection: ", err)
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

func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}
	runtime.GOMAXPROCS(*s.Options.Threads + 2) // thread count + main + web server
}

func (s *Session) AddFinding(finding *signatures.Finding) {
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

func ValidateNewSession(session *Session) error {
	// var err error

	// if session.Options, err = ParseOptions(); err != nil {
	//   return err
	// }

	if *session.Options.Save != "" && filehandler.FileExists(*session.Options.Save) {
		return errors.New(fmt.Sprintf("File: %s already exists.", *session.Options.Save))
	}

	if *session.Options.Load != "" {
		if !filehandler.FileExists(*session.Options.Load) {
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
