package gitprovider

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	Client *github.Client
	AdditionalParams map[string]string
	Token string
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
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	g.Client = github.NewClient(tc)

	return nil
}

func (g *GithubProvider) GetRepository(opt map[string]string) (*Repository, error) {
	owner, exists := opt["owner"]
	if !exists {
		return nil, errors.New("owner option must exist in map")
	}

	repo, exists := opt["repo"]
	if !exists {
		return nil, errors.New("repo option must exist in map")
	}

	r, _, err := g.Client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	return &Repository{
		ID:            r.GetID(),
		Name:          r.GetName(),
		FullName:      r.GetFullName(),
		CloneURL:      r.GetCloneURL(),
		URL:           r.GetURL(),
		DefaultBranch: r.GetDefaultBranch(),
		Description:   r.GetDescription(),
		Homepage:      r.GetHomepage(),
		Owner:         r.GetOwner().GetName(),
	}, nil
}

func (g *GithubProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}

func (g *GithubProvider) Name() string {
	return "github"
}
