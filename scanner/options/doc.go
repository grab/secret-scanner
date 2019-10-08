package options

import "errors"

var (
	// ErrRepoOptionConflict defines repo option conflict error
	ErrRepoOptionConflict = errors.New("error: options repo-id and repo-list are mutually exclusive, please provide either one")
	// ErrInvalidGitProvider defines missing git provider value error
	ErrInvalidGitProvider = errors.New("error: option git empty or invalid")
)
