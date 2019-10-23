/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/bitbucket"
)

// OAuth2Config contains OAuth2 configs for Bitbucket
type OAuth2Config struct {
	oauth2.Config
	BaseURL  string
	Username string
	Password string
}

// Bitbucket service client
type Bitbucket struct {
	Client *http.Client
	config *OAuth2Config
	token  *oauth2.Token
}

// UserRepository fetches a user's repository
func (bb *Bitbucket) UserRepository(userSlug, repoSlug string) (*Repository, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", bb.config.BaseURL, path.Join("repositories", userSlug, repoSlug)), nil)
	if err != nil {
		return nil, err
	}
	if bb.token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer ", bb.token.AccessToken))
	}

	resp, err := bb.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrResponseNotOK
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	repo := &Repository{}
	err = json.Unmarshal(respBody, repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// NewClient generates a new Bitbucket service client
func NewClient(baseURL string, client *http.Client) (*Bitbucket, error) {
	apiURL := baseURL
	if apiURL == "" {
		apiURL = DefaultBaseURL
	}
	return &Bitbucket{
		Client: client,
		config: &OAuth2Config{
			BaseURL:  apiURL,
			Username: "",
			Password: "",
		},
		token: nil,
	}, nil
}

// NewOauth2Client generates a new Bitbucket service client with OAuth2 cred.
func NewOauth2Client(OAuthKey, OAuthSecret, username, password string, client *http.Client, endpoint *oauth2.Endpoint) (*Bitbucket, error) {
	oauth2Endpoint := oauth2.Endpoint{
		AuthURL:  bitbucket.Endpoint.AuthURL,
		TokenURL: bitbucket.Endpoint.TokenURL,
	}
	if endpoint != nil {
		if oauth2Endpoint.TokenURL != "" {
			oauth2Endpoint.TokenURL = endpoint.TokenURL
		}
		if oauth2Endpoint.AuthURL != "" {
			oauth2Endpoint.AuthURL = endpoint.AuthURL
		}
	}

	config := &OAuth2Config{
		Config: oauth2.Config{
			ClientID:     OAuthKey,
			ClientSecret: OAuthSecret,
			Endpoint:     oauth2Endpoint,
		},
		BaseURL:  DefaultBaseURL,
		Username: username,
		Password: password,
	}

	token, err := newToken(client, config)
	if err != nil {
		return nil, err
	}

	bb := &Bitbucket{
		Client: client,
		config: config,
		token:  token,
	}

	return bb, nil
}

func newToken(client *http.Client, config *OAuth2Config) (*oauth2.Token, error) {
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", config.Username)
	form.Set("password", config.Password)

	req, err := http.NewRequest(http.MethodPost, config.Config.Endpoint.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Config.ClientID, config.Config.ClientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrResponseNotOK
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenDTO := &AccessTokenResponse{}
	err = json.Unmarshal(respBody, tokenDTO)
	if err != nil {
		return nil, err
	}

	expiryInNanoSec := tokenDTO.Expiry * int64(time.Second)
	expiry := time.Now().UTC().Add(time.Duration(expiryInNanoSec))

	return &oauth2.Token{
		AccessToken:  tokenDTO.AccessToken,
		TokenType:    tokenDTO.TokenType,
		RefreshToken: tokenDTO.RefreshToken,
		Expiry:       expiry,
	}, nil
}
