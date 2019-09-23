package bitbucket

import "errors"

const(
	DefaultBaseURL = "https://api.bitbucket.org/2.0"
)

var (
	ErrResponseNotOK = errors.New("response is not 200")
)
