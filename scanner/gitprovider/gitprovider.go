package gitprovider

// GitProvider defines interface for interacting with remote Git services
type GitProvider interface {
	Initialize(baseURL, token string, additionalParams map[string]string) error
	ValidateAdditionalParams(additionalParams map[string]string) bool
	GetRepository(opt map[string]string) (*Repository, error)
	Name() string
}
