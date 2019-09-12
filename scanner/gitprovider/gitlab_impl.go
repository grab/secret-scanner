package gitprovider

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabProvider struct {
	Client *gitlab.Client
	AdditionalParams map[string]string
	token string
}

func (g *GitlabProvider) Initialize(baseURL, token string, additionalParams map[string]string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	if !g.ValidateAdditionalParams(additionalParams) {
		return ErrInvalidAdditionalParams
	}

	g.token = token
	g.AdditionalParams = additionalParams
	g.Client = gitlab.NewClient(nil, token)
	err := g.Client.SetBaseURL(baseURL)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProvider) GetRepository(id string) (*Repository, error) {
	proj, _, err := g.Client.Projects.GetProject(id, nil)
	if err != nil {
		return nil, err
	}
	repo := &Repository{
		ID:            int64(proj.ID),
		Name:          proj.Name,
		FullName:      proj.Name,
		CloneURL:      proj.SSHURLToRepo,
		URL:           proj.WebURL,
		DefaultBranch: proj.DefaultBranch,
		Description:   proj.Description,
		Homepage:      proj.WebURL,
		Owner:         "",
	}
	return repo, nil
}

func (g *GitlabProvider) ValidateAdditionalParams(additionalParams map[string]string) bool {
	return true
}
