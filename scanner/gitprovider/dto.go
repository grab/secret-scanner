package gitprovider

// Repository is a universal struct for holding repo info fields
type Repository struct {
	Owner         string
	ID            string
	Name          string
	FullName      string
	CloneURL      string
	URL           string
	DefaultBranch string
	Description   string
	Homepage      string
}
