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
