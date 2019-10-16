/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package bitbucket

// AccessTokenResponse fields
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Expiry       int64  `json:"expires_in,omitempty"`
}

// Repository fields
type Repository struct {
	SCM         string           `json:"scm"`
	Website     string           `json:"website"`
	HasWiki     bool             `json:"has_wiki"`
	UUID        string           `json:"uuid"`
	Links       *RepositoryLinks `json:"links"`
	ForkPolicy  string           `json:"fork_policy"`
	Name        string           `json:"name"`
	Project     *Project         `json:"project"`
	Language    string           `json:"language"`
	CreatedOn   string           `json:"created_on"`
	MainBranch  *BranchInfo      `json:"mainbranch"`
	FullName    string           `json:"full_name"`
	HasIssues   bool             `json:"has_issues"`
	Owner       *Owner           `json:"owner"`
	UpdatedOn   string           `json:"updated_on"`
	Size        int64            `json:"size"`
	Type        string           `json:"type"`
	Slug        string           `json:"slug"`
	IsPrivate   bool             `json:"is_private"`
	Description string           `json:"description"`
}

// RepositoryLinks fields
type RepositoryLinks struct {
	Watchers     *Link   `json:"watchers"`
	Branches     *Link   `json:"branches"`
	Tags         *Link   `json:"tags"`
	Commits      *Link   `json:"commits"`
	Clone        []*Link `json:"clone"`
	Self         *Link   `json:"self"`
	Source       *Link   `json:"source"`
	HTML         *Link   `json:"html"`
	Avatar       *Link   `json:"avatar"`
	Hooks        *Link   `json:"hooks"`
	Forks        *Link   `json:"forks"`
	Downloads    *Link   `json:"downloads"`
	PullRequests *Link   `json:"pullrequests"`
}

// Link fields
type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// Project fields
type Project struct {
	Key   string        `json:"key"`
	Type  string        `json:"type"`
	UUID  string        `json:"uuid"`
	Links *ProjectLinks `json:"links"`
	Name  string        `json:"name"`
}

// ProjectLinks fields
type ProjectLinks struct {
	Self   *Link `json:"self"`
	HTML   *Link `json:"html"`
	Avatar *Link `json:"avatar"`
}

// BranchInfo fields
type BranchInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Owner fields
type Owner struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	UUID        string `json:"uuid"`
}
