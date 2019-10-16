/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"testing"
)

func TestBitbucketProvider_Initialize(t *testing.T) {
	provider := createNewBitbucketProvider()
	err := provider.Initialize(server.URL, "my-token", nil)
	if err != nil {
		t.Errorf("Want no err, got err: %v", err)
		return
	}
	if provider.Client == nil {
		t.Errorf("Want client, got nil")
	}
}

func TestBitbucketProvider_GetRepository(t *testing.T) {
	provider := createNewBitbucketProvider()
	opt := map[string]string{}
	err := provider.Initialize(server.URL+"/bitbucket", "my-token", nil)
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
	if repo.Name != "mama" {
		t.Errorf("Want mama, got %v", repo.Name)
		return
	}
}

func TestBitbucketProvider_ValidateAdditionalParams(t *testing.T) {
	provider := createNewBitbucketProvider()
	if !provider.ValidateAdditionalParams(map[string]string{}) {
		t.Errorf("Want true, false")
	}
}

func TestBitbucketProvider_Name(t *testing.T) {
	provider := createNewBitbucketProvider()
	if provider.Name() != BitbucketName {
		t.Errorf("Want %v, got %v", BitbucketName, provider.Name())
	}
}

func createNewBitbucketProvider() *BitbucketProvider {
	return &BitbucketProvider{
		Client:           nil,
		AdditionalParams: nil,
		Token:            "",
	}
}
