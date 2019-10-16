/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import "errors"

const (
	// GitlabName ...
	GitlabName = "gitlab"
	// GithubName ...
	GithubName = "github"
	// BitbucketName ...
	BitbucketName = "bitbucket"
)

var (
	// ErrInvalidAdditionalParams ...
	ErrInvalidAdditionalParams = errors.New("invalid additional params")
)
