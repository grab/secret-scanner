/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"testing"
)

func TestGithubProvider_Initialize(t *testing.T) {
	provider := createNewGithubProvider()
	err := provider.Initialize(server.URL, "my-token", nil)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if provider.Client == nil {
		t.Errorf("Want client, got nil")
	}
}

func TestGithubProvider_GetRepository(t *testing.T) {
	provider := createNewGithubProvider()
	opt := map[string]string{}
	err := provider.Initialize(server.URL+"/github/", "my-token", nil)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}

	_, err = provider.GetRepository(opt)
	if err == nil {
		t.Errorf("Want err, got no err")
		return
	}

	opt["owner"] = "my-owner"
	opt["repo"] = "repo"
	repo, err := provider.GetRepository(opt)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if repo.Name != "jquery" {
		t.Errorf("Want jquery, got %v", repo.Name)
		return
	}
}

func TestGithubProvider_ValidateAdditionalParams(t *testing.T) {
	provider := createNewGithubProvider()
	if !provider.ValidateAdditionalParams(map[string]string{}) {
		t.Errorf("Want true, false")
	}
}

func TestGithubProvider_Name(t *testing.T) {
	provider := createNewGithubProvider()
	if provider.Name() != GithubName {
		t.Errorf("Want %v, got %v", GithubName, provider.Name())
	}
}

func createNewGithubProvider() *GithubProvider {
	return &GithubProvider{
		Client:           nil,
		AdditionalParams: nil,
		Token:            "",
	}
}
