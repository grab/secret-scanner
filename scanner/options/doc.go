package options

import "errors"

var (
	ErrRepoOptionConflict = errors.New("options repo-id and repo-list are mutually exclusive, please provide either one")
)
