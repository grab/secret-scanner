/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package options

import (
	"flag"
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
	ScanTarget:  flag.String("scan-target", "", "Sub-directory within the repository to scan"),
	Repos:       flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
	LocalPath:   flag.String("git-scan-path", "", "Specify the local path to scan"),
	UI:          flag.Bool("ui", true, "Serves up local UI for scan results if true, defaults to true"),
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
