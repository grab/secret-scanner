package gitprovider

import (
	"errors"
	"net/http"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/external/remotegit/bitbucket"
)

// BitbucketProvider holds Bitbucket client fields
type BitbucketProvider struct {
	Client           *bitbucket.Bitbucket
	AdditionalParams map[string]string
	Token            string
}

// Initialize creates and assigns new client
func (g *BitbucketProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	bb, err := bitbucket.NewClient(baseURL, http.DefaultClient)
	if err != nil {
		return err
	}
	g.Client = bb
	g.AdditionalParams = additionalParams

	return nil
}

// GetRepository gets repo info
func (g *BitbucketProvider) GetRepository(opt map[string]string) (*Repository, error) {
	username, exists := opt["owner"]
	if !exists {
		return nil, errors.New("username option must exist in map")
	}

	repoSlug, exists := opt["repo"]
	if !exists {
		return nil, errors.New("repoSlug option must exist in map")
	}

	repo, err := g.Client.UserRepository(username, repoSlug, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Owner:         repo.Owner.Username,
		ID:            repo.UUID,
		Name:          repo.Name,
		FullName:      repo.FullName,
		CloneURL:      repo.Links.Clone[0].Href,
		URL:           repo.Links.Self.Href,
		DefaultBranch: repo.MainBranch.Name,
		Description:   repo.Description,
		Homepage:      repo.Links.HTML.Href,
	}, nil
}

// ValidateAdditionalParams validates additional params
func (g *BitbucketProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}

// Name returns the provider name
func (g *BitbucketProvider) Name() string {
	return BitbucketName
}
