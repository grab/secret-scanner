package gitprovider

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	Client *github.Client
	Token string
	AdditionalParams map[string]string
}

func (g *GithubProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	if !g.ValidateAdditionalParams(additionalParams) {
		return ErrInvalidAdditionalParams
	}

	g.Token = token
	g.AdditionalParams = additionalParams
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	g.Client = github.NewClient(tc)

	return nil
}

func (g *GithubProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}
