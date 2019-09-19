package gitprovider

import "errors"

const (
	GitlabName = "gitlab"
	GithubName = "github"
)

var (
	ErrEmptyToken = errors.New("token must not be empty")
	ErrInvalidAdditionalParams = errors.New("invalid additional params")
)
