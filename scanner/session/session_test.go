/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package session

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/grab/secret-scanner/scanner/gitprovider"

	"github.com/grab/secret-scanner/scanner/findings"

	"github.com/grab/secret-scanner/scanner/options"
)

var defaultOptions = options.Options{
	CommitDepth: flag.Int("commit-depth", 500, "Number of repository commits to process"),
	Threads:     flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
	Report:      flag.String("save", "", "Save session to file"),
	Load:        flag.String("load", "", "Load session file"),
	Silent:      flag.Bool("silent", false, "Suppress all output except for errors"),
	Debug:       flag.Bool("debug", false, "Print debugging information"),
	GitProvider: flag.String("git", "", "Specify type of git provider (Eg. github, gitlab, bitbucket)"),
	BaseURL:     flag.String("baseurl", "", "Specify Git provider base URL"),
	Token:       flag.String("token", "", "Specify Git provider token"),
	EnvFilePath: flag.String("env", "", ".env file path containing Git provider base URLs and tokens"),
	ScanTarget:  flag.String("scan-target", "", "Sub-directory within the repository to scan"),
	Repos:       flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
	LocalPath:   flag.String("git-scan-path", "", "Specify the local path to scan"),
	UI:          flag.Bool("ui", true, "Serves up local UI for scan results if true, defaults to true"),
}

func TestSession_Initialize(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	if *sess.Options.CommitDepth != 500 {
		t.Errorf("Want 1, got %v", *sess.Options.CommitDepth)
	}
	if sess.Out == nil {
		t.Errorf("Want Logger, got nil")
	}
	if sess.Stats == nil {
		t.Errorf("Want Stats, got nil")
	}
	if sess.StateStore == nil {
		t.Errorf("Want StateStore, got nil")
	}
	sess.End()
}

func TestSession_End(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	sess.End()

	nilTime := time.Time{}
	if sess.Stats.FinishedAt == nilTime {
		t.Errorf("Stats field FinishAt should not be nil time")
	}
	if sess.Stats.Status != StatusFinished {
		t.Errorf("Want %v, got %v", StatusFinished, sess.Stats.Status)
	}
}

func TestSession_InitLogger(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	if sess.Out == nil {
		t.Errorf("Want Logger, got nil")
	}
	sess.End()
}

func TestSession_InitStats(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	if sess.Stats == nil {
		t.Errorf("Want Stats, got nil")
	}
	sess.End()
}

func TestSession_AddFinding(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	sess.AddFinding(&findings.Finding{})
	if len(sess.Findings) != 1 {
		t.Errorf("Want 1, got %v", len(sess.Findings))
	}
	sess.End()
}

func TestSession_AddRepository(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)
	sess.AddRepository(&gitprovider.Repository{})
	if len(sess.Repositories) != 1 {
		t.Errorf("Want 1, got %v", len(sess.Repositories))
	}
	sess.End()
}

func TestSession_SaveToFile(t *testing.T) {
	sess := createNewSession()
	sess.Initialize(defaultOptions)

	tempDir, err := ioutil.TempDir("", "ss-test-")
	if err != nil {
		t.Errorf("Cannot create temp. dir.: %v", err)
		return
	}

	filepath := path.Join(tempDir, "ss-test.json")
	absPath, err := sess.SaveToFile(filepath)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if absPath != filepath {
		t.Errorf("Want %v, got %v", filepath, absPath)
		return
	}

	_, err = os.Stat(filepath)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
	}

	_ = os.RemoveAll(tempDir)

	sess.End()
}

func createNewSession() *Session {
	return &Session{
		Mutex:        sync.Mutex{},
		Options:      options.Options{},
		Out:          nil,
		Stats:        nil,
		Findings:     nil,
		Repositories: nil,
		Signatures:   nil,
		StateStore:   nil,
	}
}
