package gitprovider

type GitProvider interface {
	Initialize(baseURL, token string, additionalParams map[string]string) error
	ValidateAdditionalParams(additionalParams map[string]string) bool
	GetRepository(id string) (*Repository, error)
}
