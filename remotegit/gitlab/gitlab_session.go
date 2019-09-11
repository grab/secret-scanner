package gitlab

import (
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/logic/scan"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/remotegit"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
	"os"

	"github.com/xanzy/go-gitlab"
)

type GitlabSession struct {
	*session.Session
	GitlabAccessToken string         `json:"-"`
	GitlabClient      *gitlab.Client `json:"-"`
	GitlabRepos       []*remotegit.Repository
}

func (s *GitlabSession) Start() {
	s.Session.Initialize("gitlab")
	s.InitGitlabAccessToken()
	s.InitGitlabClient()
}

func (s *GitlabSession) Finish() {
	s.Session.End()
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

func (s *GitlabSession) AddGitlabRepository(repository *remotegit.Repository) {
	s.Lock()
	defer s.Unlock()
	for _, r := range s.GitlabRepos {
		if repository.ID == r.ID {
			return
		}
	}
	s.GitlabRepos = append(s.GitlabRepos, repository)
}

func NewGitlabSession(options scan.Options) (*GitlabSession, error) {
	var err error
	var gitlabRepos []*remotegit.Repository
	sess := GitlabSession{&session.Session{}, "", nil, gitlabRepos}
	sess.Options = options
	err = session.ValidateNewSession(sess.Session)
	if err != nil {
		return nil, err
	}
	sess.Start()
	return &sess, nil
}
