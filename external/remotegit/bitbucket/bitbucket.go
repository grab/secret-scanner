package bitbucket

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/bitbucket"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type OAuthConfig struct {
	oauth2.Config
	BaseURL string
	Username string
	Password string
}

type Bitbucket struct {
	Client *http.Client
	config *OAuthConfig
	token *oauth2.Token
}

func (bb *Bitbucket) UserRepository(userSlug, repoSlug string, client *http.Client) (*Repository, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", bb.config.BaseURL, path.Join("repositories", userSlug, repoSlug)), nil)
	if err != nil {
		return nil, err
	}
	if bb.token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer ", bb.token.AccessToken))
	}

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

	repo := &Repository{}
	err = json.Unmarshal(respBody, repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func NewClient(client *http.Client) (*Bitbucket, error) {
	return &Bitbucket{
		Client: client,
		config: &OAuthConfig{
			BaseURL: DefaultBaseURL,
			Username: "",
			Password: "",
		},
		token:  nil,
	}, nil
}

func NewOauthClient(OAuthKey, OAuthSecret, username, password string, client *http.Client) (*Bitbucket, error) {
	config := &OAuthConfig{
		Config: oauth2.Config{
			ClientID: OAuthKey,
			ClientSecret: OAuthSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   bitbucket.Endpoint.AuthURL,
				TokenURL:  bitbucket.Endpoint.TokenURL,
			},
		},
		BaseURL: DefaultBaseURL,
		Username: username,
		Password: password,
	}

	token, err := newToken(client, config)
	if err != nil {
		return nil, err
	}

	bb := &Bitbucket{
		Client:  client,
		config:  config,
		token:   token,
	}

	return bb, nil
}

func newToken(client *http.Client, config *OAuthConfig) (*oauth2.Token, error) {
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
		AccessToken: tokenDTO.AccessToken,
		TokenType: tokenDTO.TokenType,
		RefreshToken: tokenDTO.RefreshToken,
		Expiry: expiry,
	}, nil
}
