/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/grab/secret-scanner/scanner"
	"github.com/grab/secret-scanner/scanner/gitprovider"
	"github.com/grab/secret-scanner/scanner/options"
	"github.com/grab/secret-scanner/scanner/session"
)

func main() {
	// Parse CLI options
	opt, err := options.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Load env file
	loadEnv(*opt.EnvFilePath)

	var gitProvider gitprovider.GitProvider
	additionalParams := map[string]string{}

	// Set Git provider
	switch *opt.GitProvider {
	case gitprovider.GithubName:
		gitProvider = &gitprovider.GithubProvider{}
		if *opt.BaseURL == "" {
			*opt.BaseURL = os.Getenv(gitprovider.GithubParamBaseURL)
		}
		if *opt.Token == "" {
			*opt.Token = os.Getenv(gitprovider.GithubParamToken)
		}
	case gitprovider.GitlabName:
		gitProvider = &gitprovider.GitlabProvider{}
		if *opt.BaseURL == "" {
			*opt.BaseURL = os.Getenv(gitprovider.GitlabParamBaseURL)
		}
		if *opt.Token == "" {
			*opt.Token = os.Getenv(gitprovider.GitlabParamToken)
		}
	case gitprovider.BitbucketName:
		gitProvider = &gitprovider.BitbucketProvider{}
		if *opt.BaseURL == "" {
			*opt.BaseURL = os.Getenv(gitprovider.BitbucketParamBaseURL)
		}
		additionalParams[gitprovider.BitbucketParamClientID] = os.Getenv(gitprovider.BitbucketParamClientID)
		additionalParams[gitprovider.BitbucketParamClientSecret] = os.Getenv(gitprovider.BitbucketParamClientSecret)
		additionalParams[gitprovider.BitbucketParamUsername] = os.Getenv(gitprovider.BitbucketParamUsername)
		additionalParams[gitprovider.BitbucketParamPassword] = os.Getenv(gitprovider.BitbucketParamPassword)
	default:
		fmt.Println("error: invalid Git provider type (Currently supports github, gitlab, bitbucket)")
		os.Exit(1)
	}

	// Initialize Git provider
	err = gitProvider.Initialize(*opt.BaseURL, *opt.Token, additionalParams)
	if err != nil {
		fmt.Println(errors.New(fmt.Sprintf("unable to initialise %s provider", *opt.GitProvider)))
		os.Exit(1)
	}

	// Initialize new scan session
	sess := &session.Session{}
	sess.Initialize(opt)
	sess.Out.Important("%s Scanning Started at %s\n", strings.Title(*opt.GitProvider), sess.Stats.StartedAt.Format(time.RFC3339))
	sess.Out.Important("Loaded %d signatures\n", len(sess.Signatures))

	if sess.Stats.Status == "finished" {
		sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
		return
	}

	// Scan
	scanner.Scan(sess, gitProvider)
	sess.Out.Important("Gitlab Scanning Finished at %s\n", sess.Stats.FinishedAt.Format(time.RFC3339))

	if *sess.Options.Report != "" {
		absPath, err := sess.SaveToFile(*sess.Options.Report)
		if err != nil {
			sess.Out.Error("Error saving session to %s: %s\n", *sess.Options.Report, err)
		}
		sess.Out.Important("Saved session to: %s\n\n", absPath)
	}

	sess.Stats.PrintStats(sess.Out)
}

func loadEnv(envPath string) {
	if envPath != "" {
		err := godotenv.Load(envPath)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v, seaching in work directory instead", err))
		} else {
			return
		}
	}

	currentWD, err := os.Getwd()
	if err != nil {
		fmt.Println(fmt.Sprintf("error getting work directory"))
		return
	}

	envPath = path.Join(currentWD, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println(err)
	}
}
