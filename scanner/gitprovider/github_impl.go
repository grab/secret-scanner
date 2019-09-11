package gitprovider

import "github.com/google/go-github/github"

type GithubProvider struct {
	Client *github.Client
	Token string
	AdditionalParams map[string]string
}

func (g *GithubProvider) ValidateAdditionalParams() bool {
	return true
}
