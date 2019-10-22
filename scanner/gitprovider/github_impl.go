/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"

	//"golang.org/x/oauth2"
	"strconv"
)

// GithubProvider holds Github client fields
type GithubProvider struct {
	Client           *github.Client
	AdditionalParams map[string]string
	Token            string
}

// Initialize creates and assigns new client
func (g *GithubProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	if !g.ValidateAdditionalParams(additionalParams) {
		return ErrInvalidAdditionalParams
	}

	var client *http.Client
	g.AdditionalParams = additionalParams

	if token != "" {
		g.Token = token
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		client = oauth2.NewClient(context.Background(), ts)
	}

	g.Client = github.NewClient(client)

	// change client's base URL if needed
	if baseURL != "" {
		apiURL, _ := url.Parse(baseURL)
		g.Client.BaseURL = apiURL
	}

	return nil
}

// GetRepository gets repo info
func (g *GithubProvider) GetRepository(opt map[string]string) (*Repository, error) {
	owner, exists := opt["owner"]
	if !exists {
		return nil, errors.New("owner option must exist in map")
	}

	repo, exists := opt["repo"]
	if !exists {
		return nil, errors.New("repo option must exist in map")
	}

	r, _, err := g.Client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	return &Repository{
		ID:            strconv.Itoa(int(r.GetID())),
		Name:          r.GetName(),
		FullName:      r.GetFullName(),
		CloneURL:      r.GetCloneURL(),
		URL:           r.GetURL(),
		DefaultBranch: r.GetDefaultBranch(),
		Description:   r.GetDescription(),
		Homepage:      r.GetHomepage(),
		Owner:         r.GetOwner().GetName(),
	}, nil
}

// GetAdditionalParams validates additional params
func (g *GithubProvider) GetAdditionalParam(key string) string {
	val, exists := g.AdditionalParams[key]
	if !exists {
		return ""
	}
	return val
}

// ValidateAdditionalParams validates additional params
func (g *GithubProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}

// Name returns the provider name
func (g *GithubProvider) Name() string {
	return GithubName
}
