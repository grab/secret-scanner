/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"errors"
	"net/http"

	"github.com/grab/secret-scanner/external/remotegit/bitbucket"
)

// BitbucketProvider holds Bitbucket client fields
type BitbucketProvider struct {
	Client           *bitbucket.Bitbucket
	AdditionalParams map[string]string
	Token            string
}

// Initialize creates and assigns new client
func (g *BitbucketProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	if !g.ValidateAdditionalParams(additionalParams) {
		return ErrInvalidAdditionalParams
	}

	var bb *bitbucket.Bitbucket
	var err error
	g.AdditionalParams = additionalParams

	if g.AdditionalParams[BitbucketParamClientID] != "" &&
		g.AdditionalParams[BitbucketParamClientSecret] != "" &&
		g.AdditionalParams[BitbucketParamUsername] != "" &&
		g.AdditionalParams[BitbucketParamPassword] != "" {

		bb, err = bitbucket.NewOauth2Client(
			g.AdditionalParams[BitbucketParamClientID],
			g.AdditionalParams[BitbucketParamClientSecret],
			g.AdditionalParams[BitbucketParamUsername],
			g.AdditionalParams[BitbucketParamPassword],
			http.DefaultClient,
			nil)
		if err != nil {
			return err
		}

		g.Client = bb

		return nil
	}

	bb, err = bitbucket.NewClient(baseURL, http.DefaultClient)
	if err != nil {
		return err
	}
	g.Client = bb

	return nil
}

// GetRepository gets repo info
func (g *BitbucketProvider) GetRepository(opt map[string]string) (*Repository, error) {
	username, exists := opt["owner"]
	if !exists {
		return nil, errors.New("username option must exist in map")
	}

	repoSlug, exists := opt["repo"]
	if !exists {
		return nil, errors.New("repoSlug option must exist in map")
	}

	repo, err := g.Client.UserRepository(username, repoSlug)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Owner:         repo.Owner.Username,
		ID:            repo.UUID,
		Name:          repo.Name,
		FullName:      repo.FullName,
		CloneURL:      repo.Links.Clone[0].Href,
		URL:           repo.Links.Self.Href,
		DefaultBranch: repo.MainBranch.Name,
		Description:   repo.Description,
		Homepage:      repo.Links.HTML.Href,
	}, nil
}

// GetAdditionalParams validates additional params
func (g *BitbucketProvider) GetAdditionalParam(key string) string {
	val, exists := g.AdditionalParams[key]
	if !exists {
		return ""
	}
	return val
}

// ValidateAdditionalParams validates additional params
func (g *BitbucketProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}

// Name returns the provider name
func (g *BitbucketProvider) Name() string {
	return BitbucketName
}
