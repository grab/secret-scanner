package gitprovider

type Repository struct {
	Owner         string
	ID            int64
	Name          string
	FullName      string
	CloneURL      string
	URL           string
	DefaultBranch string
	Description   string
	Homepage      string
}
