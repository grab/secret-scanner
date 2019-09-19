package options

import "errors"

var (
	ErrRepoOptionConflict = errors.New("error: options repo-id and repo-list are mutually exclusive, please provide either one")
	ErrInvalidGitProvider = errors.New("error: option vcs empty or invalid")
)
