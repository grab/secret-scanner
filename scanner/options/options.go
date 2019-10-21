/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package options

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/filehandler"

	"github.com/mitchellh/go-homedir"

	"github.com/joho/godotenv"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
)

// Options ...
type Options struct {
	CommitDepth      *int
	Threads          *int
	Report           *string `json:"-"`
	Load             *string `json:"-"`
	Silent           *bool
	Debug            *bool
	SkipTestContexts *bool
	NoHistory        *bool
	LogSecret        *bool

	GitProvider *string
	BaseURL     *string
	Token       *string
	//ClientID           *string
	//ClientSecret       *string
	//UserID             *string
	//UserPW             *string
	EnvFilePath          *string
	HistoryStoreFilePath *string
	RepoID               *string
	ScanTarget           *string
	Repos                *string
	GitScanPath          *string
	UI                   *bool
}

// ValidateOptions validates given options
func (o Options) ValidateOptions() (bool, error) {
	if *o.RepoID != "" && *o.Repos != "" {
		return false, ErrRepoOptionConflict
	}
	if *o.EnvFilePath != "" {
		if _, err := os.Stat(*o.EnvFilePath); os.IsNotExist(err) {
			return false, err
		}
	}

	// Load env file if present
	if *o.EnvFilePath != "" {
		err := godotenv.Load(*o.EnvFilePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("error: unable to load .env file path %s: %v", *o.EnvFilePath, err))
			os.Exit(1)
		}
	}

	if *o.GitProvider != gitprovider.GithubName && *o.GitProvider != gitprovider.GitlabName && *o.GitProvider != gitprovider.BitbucketName {
		return false, ErrInvalidGitProvider
	}

	switch *o.GitProvider {
	case gitprovider.GithubName:
		return o.ValidateGithubOptions(), nil
	case gitprovider.GitlabName:
		return o.ValidateGitlabOptions(), nil
	case gitprovider.BitbucketName:
		return o.ValidateBitbucketOptions(), nil
	default:
		return false, ErrInvalidGitProvider
	}
}

// ValidateGithubOptions validates Github options,
// applied when GitProvider == github
func (o Options) ValidateGithubOptions() bool {
	return true
}

// ValidateGitlabOptions validates Gitlab options
// applied when GitProvider == gitlab
func (o Options) ValidateGitlabOptions() bool {
	if o.BaseURL == nil {
		return false
	}

	baseURL := *o.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("GITLAB_BASE_URL")
		*o.BaseURL = baseURL
	}
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return false
	}
	return o.ValidateHasToken("GITLAB_TOKEN")
}

// ValidateBitbucketOptions validates Bitbucket options
// applied when GitProvider == bitbucket
func (o Options) ValidateBitbucketOptions() bool {
	return true
}

// ValidateHasToken validates that token is not empty
func (o *Options) ValidateHasToken(key string) bool {
	if *o.Token == "" {
		if os.Getenv(key) == "" {
			return false
		}
		//token := os.Getenv(key)
		*o.Token = os.Getenv(key)
	}
	return true
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
		CommitDepth:      flag.Int("commit-depth", 500, "Number of repository commits to process"),
		Threads:          flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Report:           flag.String("report", "", "Save session to file"),
		Load:             flag.String("load", "", "Load session file"),
		Silent:           flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:            flag.Bool("debug", false, "Print debugging information"),
		SkipTestContexts: flag.Bool("skip-tests", false, "Skips possible test contexts"),
		NoHistory:        flag.Bool("no-history", true, "If no-history is on, every scan will be treated as a brand new scan."),
		LogSecret:        flag.Bool("log-secret", true, "If true, the matched secret will be included in results save file"),

		GitProvider: flag.String("git", "", "Specify type of git provider (Eg. github, gitlab, bitbucket)"),
		BaseURL:     flag.String("baseurl", "", "Specify Git provider base URL"),
		Token:       flag.String("token", "", "Specify Git provider token"),
		//ClientID:           flag.String("oauth-id", "", "Specify Bitbucket Oauth2 client ID"),
		//ClientSecret:       flag.String("oauth-secret", "", "Specify Bitbucket Oauth2 client secret"),
		//UserID:             flag.String("user-id", "", "Specify Bitbucket username"),
		//UserPW:             flag.String("user-pw", "", "Specify Bitbucket password"),
		EnvFilePath:          flag.String("env", "", ".env file path containing Git provider base URLs and tokens"),
		HistoryStoreFilePath: flag.String("history", "", "File path to store scan histories"),
		RepoID:               flag.String("repo-id", "", "Scan the repository with this ID"),
		ScanTarget:           flag.String("scan-target", "", "Sub-directory within the repository to scan"),
		Repos:                flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
		GitScanPath:          flag.String("git-scan-path", "", "Specify the local path to scan"),
		UI:                   flag.Bool("ui", false, "Serves up local UI for scan results if true, defaults to true"),
	}

	flag.Parse()

	return options, nil
}

func LoadDefaultConfig(filename string) (*Options, error) {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	if filename == "" {
		filename = DefaultConfigFilename
	}

	cfgDirPath := path.Join(userHomeDir, DefaultLocation)
	cfgFilePath := path.Join(cfgDirPath, filename)

	if !filehandler.FileExists(cfgFilePath) {
		return CreateDefaultConfig(), nil
	}
	return nil, nil
}

func CreateDefaultConfig() *Options {
	return &Options{
		CommitDepth:          nil,
		Threads:              nil,
		Report:               nil,
		Load:                 nil,
		Silent:               nil,
		Debug:                nil,
		SkipTestContexts:     nil,
		NoHistory:            nil,
		LogSecret:            nil,
		GitProvider:          nil,
		BaseURL:              nil,
		Token:                nil,
		EnvFilePath:          nil,
		HistoryStoreFilePath: nil,
		RepoID:               nil,
		ScanTarget:           nil,
		Repos:                nil,
		GitScanPath:          nil,
		UI:                   nil,
	}
}
