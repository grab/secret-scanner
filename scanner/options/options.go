/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package options

import (
	"flag"
	"strings"
)

// Options ...
type Options struct {
	BaseURL          *string `json:"base_url"`
	CommitDepth      *int    `json:"commit_depth"`
	Debug            *bool   `json:"debug"`
	EnvFilePath      *string `json:"env_file_path"`
	GitProvider      *string `json:"git_provider"`
	Load             *string `json:"-"`
	LocalPath        *string `json:"local_path"`
	LogSecret        *bool   `json:"log_secret"`
	Report           *string `json:"-"`
	Repos            *string `json:"repos"`
	ScanTarget       *string `json:"scan_target"`
	Silent           *bool   `json:"silent"`
	SkipTestContexts *bool   `json:"skip_test_contexts"`
	State            *bool   `json:"state"`
	Threads          *int    `json:"threads"`
	Token            *string `json:"token"`
	UI               *bool   `json:"ui"`
	UIHost           *string `json:"ui_host"`
	UIPort           *string `json:"ui_port"`
}

// ParseScanTargets splits string of targets by comma
func (o Options) ParseScanTargets() []string {
	targets := strings.Split(*o.ScanTarget, ",")
	for i, t := range targets {
		targets[i] = strings.Trim(t, " ")
	}
	return targets
}

// Parse parses cmd params
func Parse() (Options, error) {
	options := Options{
		BaseURL:          flag.String("baseurl", "", "Specify Git provider base URL"),
		CommitDepth:      flag.Int("commit-depth", 500, "Number of repository commits to process"),
		Debug:            flag.Bool("debug", false, "Print debugging information"),
		EnvFilePath:      flag.String("env", "", ".env file path containing Git provider base URLs and tokens"),
		GitProvider:      flag.String("git", "github", "Name of git provider (Eg. github, gitlab, bitbucket)"),
		Load:             flag.String("load", "", "Load session file"),
		LocalPath:        flag.String("dir", "", "Specify the local git repo path to scan"),
		LogSecret:        flag.Bool("log-secret", true, "If true, the matched secret will be included in report file"),
		Report:           flag.String("output", "", "Save session to file"),
		Repos:            flag.String("repos", "", "Comma-separated list of repos to scan"),
		ScanTarget:       flag.String("sub-dir", "", "Sub-directory within the repository to scan"),
		Silent:           flag.Bool("quiet", false, "Suppress all output except for errors"),
		SkipTestContexts: flag.Bool("skip-tests", true, "Skips possible test contexts"),
		State:            flag.Bool("use-state", false, "If state is off, every scan will be treated as a brand new scan."),
		Threads:          flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Token:            flag.String("token", "", "Specify Git provider token"),
		//UI:               flag.Bool("ui", false, "Serves up local UI for scan results if true"),
		//UIHost:           flag.String("ui-host", "127.0.0.1", "UI server host"),
		//UIPort:           flag.String("ui-port", "8080", "UI server port"),
	}

	flag.Parse()

	return options, nil
}
