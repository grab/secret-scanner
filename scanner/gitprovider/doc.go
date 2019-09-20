package gitprovider

import "errors"

const (
	GitlabName = "gitlab"
	GithubName = "github"
	BitbucketName = "bitbucket"
)

var (
	ErrEmptyToken = errors.New("token must not be empty")
	ErrInvalidAdditionalParams = errors.New("invalid additional params")
)
