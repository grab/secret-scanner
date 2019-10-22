/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"errors"
	"strconv"

	"github.com/xanzy/go-gitlab"
)

// GitlabProvider holds Gitlab client fields
type GitlabProvider struct {
	Client           *gitlab.Client
	AdditionalParams map[string]string
	Token            string
}

// Initialize creates and assigns new client
func (g *GitlabProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	if !g.ValidateAdditionalParams(additionalParams) {
		return ErrInvalidAdditionalParams
	}

	g.Token = token
	g.AdditionalParams = additionalParams
	g.Client = gitlab.NewClient(nil, token)

	if baseURL != "" {
		err := g.Client.SetBaseURL(baseURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRepository gets repo info
func (g *GitlabProvider) GetRepository(opt map[string]string) (*Repository, error) {
	id, exists := opt["id"]
	if !exists {
		return nil, errors.New("id option does not exists in map")
	}
	proj, _, err := g.Client.Projects.GetProject(id, nil)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		ID:            strconv.Itoa(proj.ID),
		Name:          proj.Name,
		FullName:      proj.Name,
		CloneURL:      proj.SSHURLToRepo,
		URL:           proj.WebURL,
		DefaultBranch: proj.DefaultBranch,
		Description:   proj.Description,
		Homepage:      proj.WebURL,
		Owner:         "",
	}
	return repo, nil
}

// GetAdditionalParams validates additional params
func (g *GitlabProvider) GetAdditionalParam(key string) string {
	val, exists := g.AdditionalParams[key]
	if !exists {
		return ""
	}
	return val
}

// ValidateAdditionalParams validates additional params
func (g *GitlabProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}

// Name returns the provider name
func (g *GitlabProvider) Name() string {
	return GitlabName
}
