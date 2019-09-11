package gitprovider

type GitProvider interface {
	ValidateAdditionalParams() bool
}

type Providers struct {
	Gitlab *GitlabProvider
	Github *GithubProvider
}
