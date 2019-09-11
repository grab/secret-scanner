package gitlab

import (
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scan"
	"os"

	"github.com/xanzy/go-gitlab"
)

type GitlabSession struct {
	*scan.Session
	GitlabAccessToken string         `json:"-"`
	GitlabClient      *gitlab.Client `json:"-"`
	GitlabRepos       []*GitlabRepository
}

func (s *GitlabSession) Start() {
	s.Session.Start()
	s.InitGitlabAccessToken()
	s.InitGitlabClient()
}

func (s *GitlabSession) Finish() {
	s.Session.Finish()
}

func (s *GitlabSession) InitGitlabAccessToken() {
	accessToken := os.Getenv(GitlabTokenEnvVariable)
	if accessToken == "" {
		s.Out.Fatal("No Gitlab access token given. Please provide via %s environment variable.\n", GitlabTokenEnvVariable)
	}
	s.GitlabAccessToken = accessToken
}

func (s *GitlabSession) InitGitlabClient() {
	s.GitlabClient = gitlab.NewClient(nil, s.GitlabAccessToken)
	s.GitlabClient.SetBaseURL(GitlabEndpoint)
}

func (s *GitlabSession) AddGitlabRepository(repository *GitlabRepository) {
	s.Lock()
	defer s.Unlock()
	for _, r := range s.GitlabRepos {
		if *repository.ID == *r.ID {
			return
		}
	}
	s.GitlabRepos = append(s.GitlabRepos, repository)
}

func NewGitlabSession(options scan.Options) (*GitlabSession, error) {
	var err error
	var gitlabRepos []*GitlabRepository
	session := GitlabSession{&scan.Session{}, "", nil, gitlabRepos}
	session.Options = options
	err = scan.ValidateNewSession(session.Session)
	if err != nil {
		return nil, err
	}
	session.Start()
	return &session, nil
}
