package gitprovider

type GitProvider interface {
	Initialize(baseURL, token string, additionalParams map[string]string) error
	ValidateAdditionalParams(additionalParams map[string]string) bool
	GetRepository(opt map[string]string) (*Repository, error)
	Name() string
}
