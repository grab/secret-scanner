package bitbucket

import "errors"

const (
	// DefaultBaseURL defines the default Bitbucket API URL
	DefaultBaseURL = "https://api.bitbucket.org/2.0"
)

var (
	// ErrResponseNotOK defines non-200 HTTP response error
	ErrResponseNotOK = errors.New("response is not 200")
)
