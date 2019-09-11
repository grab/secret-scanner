package gitprovider

import "github.com/xanzy/go-gitlab"

type GitlabProvider struct {
	Client *gitlab.Client
	Token string
	AdditionalParams map[string]string
}

func (g *GitlabProvider) ValidateAdditionalParams() bool {
	return true
}
