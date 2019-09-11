package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/log"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/logic/scan"
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

	Options  scan.Options `json:"-"`
	Out      *log.Logger  `json:"-"`
	Stats    *stats.Stats
	Findings []*scan.Finding
	Store    *db.MysqlHandler
	//RemoteGitClients *RemoteGitClients
}

func (s *Session) Initialize(remoteGitType string) {
	//s.RemoteGitClients = &RemoteGitClients{}

	//switch remoteGitType {
	//case "gitlab":
	//	s.RemoteGitClients.Gitlab = gitlab.NewClient(nil, s.GitlabAccessToken)
	//}
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
	s.InitLogger()
	s.InitThreads()
	s.InitDB()
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

func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}
	runtime.GOMAXPROCS(*s.Options.Threads + 2) // thread count + main + web server
}

func (s *Session) AddFinding(finding *scan.Finding) {
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

func ValidateNewSession(session *Session) error {
	// var err error

	// if session.Options, err = ParseOptions(); err != nil {
	//   return err
	// }

	if *session.Options.Save != "" && scan.FileExists(*session.Options.Save) {
		return errors.New(fmt.Sprintf("File: %s already exists.", *session.Options.Save))
	}

	if *session.Options.Load != "" {
		if !scan.FileExists(*session.Options.Load) {
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
