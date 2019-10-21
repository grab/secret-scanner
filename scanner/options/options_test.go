/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package options

import (
	"flag"
	"os"
	"testing"
)

var defaultOptions = Options{
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
	RepoID:      flag.String("repo-id", "", "Scan the repository with this ID"),
	ScanTarget:  flag.String("scan-target", "", "Sub-directory within the repository to scan"),
	Repos:       flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
	GitScanPath: flag.String("git-scan-path", "", "Specify the local path to scan"),
	UI:          flag.Bool("ui", true, "Serves up local UI for scan results if true, defaults to true"),
}

func TestOptions_ValidateOptions(t *testing.T) {
	o := defaultOptions

	valid, err := o.ValidateOptions()
	if err != ErrInvalidGitProvider {
		t.Errorf("Want ErrInvalidGitProvider, no err")
	}
	if valid {
		t.Errorf("Want false, got true")
	}

	*o.GitProvider = "github"
	valid, err = o.ValidateOptions()
	if err != nil {
		t.Errorf("Want no err, got err")
	}
	if !valid {
		t.Errorf("Want true, got false")
	}

	*o.RepoID = "123"
	*o.Repos = "some/path"
	valid, err = o.ValidateOptions()
	if err != ErrRepoOptionConflict {
		t.Errorf("Want ErrRepoOptionConflict, no err")
	}
	if valid {
		t.Errorf("Want false, got true")
	}

	*o.RepoID = ""
	valid, err = o.ValidateOptions()
	if err != nil {
		t.Errorf("Want no err, got err")
	}
	if !valid {
		t.Errorf("Want true, got false")
	}
}

func TestOptions_ValidateGithubOptions(t *testing.T) {
	options := Options{}
	if !options.ValidateGithubOptions() {
		t.Errorf("Want true, got false")
	}
}

func TestOptions_ValidateGitlabOptions(t *testing.T) {
	options := Options{}
	if options.ValidateGitlabOptions() {
		t.Errorf("Want false, got true")
		return
	}

	baseURL := "https://my-gitlab.com"
	token := "my-token"
	options.BaseURL = &baseURL
	options.Token = &token

	if !options.ValidateGitlabOptions() {
		t.Errorf("Want true, got false")
		return
	}

	if *options.BaseURL != "https://my-gitlab.com" {
		t.Errorf("Want https://my-gitlab.com, got %v", *options.BaseURL)
	}

	if *options.Token != "my-token" {
		t.Errorf("Want my-token, got %v", *options.Token)
	}
}

func TestOptions_ValidateBitbucketOptions(t *testing.T) {
	options := Options{}
	if !options.ValidateBitbucketOptions() {
		t.Errorf("Want true, got false")
	}
}

func TestOptions_ValidateHasToken(t *testing.T) {
	options := defaultOptions

	if options.ValidateHasToken("SS_TEST_TOKEN") {
		t.Errorf("Want false, got true")
		return
	}

	err := os.Setenv("SS_TEST_TOKEN", "my-token-value")
	if err != nil {
		t.Errorf("Unable to set OS env. var. SS_TEST_TOKEN")
		return
	}

	if !options.ValidateHasToken("SS_TEST_TOKEN") {
		t.Errorf("Want true, got false")
		return
	}

	if *options.Token != "my-token-value" {
		t.Errorf("Want my-token-value, got %s", *options.Token)
	}
}

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
