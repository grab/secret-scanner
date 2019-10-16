/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"testing"
)

func TestGitlabProvider_Initialize(t *testing.T) {
	provider := createNewGitlabProvider()
	err := provider.Initialize(server.URL, "my-token", nil)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if provider.Client == nil {
		t.Errorf("Want client, got nil")
	}
}

func TestGitlabProvider_GetRepository(t *testing.T) {
	provider := createNewGitlabProvider()
	opt := map[string]string{}
	err := provider.Initialize(server.URL+"/gitlab", "my-token", nil)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}

	_, err = provider.GetRepository(opt)
	if err == nil {
		t.Errorf("Want err, got no err")
		return
	}

	opt["id"] = "7824084"
	repo, err := provider.GetRepository(opt)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if repo.Name != "augur" {
		t.Errorf("Want jquery, got %v", repo.Name)
		return
	}
}

func TestGitlabProvider_ValidateAdditionalParams(t *testing.T) {
	provider := createNewGitlabProvider()
	if !provider.ValidateAdditionalParams(map[string]string{}) {
		t.Errorf("Want true, false")
	}
}

func TestGitlabProvider_Name(t *testing.T) {
	provider := createNewGitlabProvider()
	if provider.Name() != GitlabName {
		t.Errorf("Want %v, got %v", GitlabName, provider.Name())
	}
}

func createNewGitlabProvider() *GitlabProvider {
	return &GitlabProvider{
		Client:           nil,
		AdditionalParams: nil,
		Token:            "",
	}
}
