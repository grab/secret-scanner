package github

import (
	"context"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/logic/scan"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubSession struct {
	*session.Session
	GithubAccessToken string         `json:"-"`
	GithubClient      *github.Client `json:"-"`
	Targets           []*GithubOwner
	Repositories      []*GithubRepository
}

func (s *GithubSession) Start() {
	s.Session.Initialize("github")
	s.InitGithubAccessToken()
	s.InitGithubClient()
}

func (s *GithubSession) Finish() {
	s.Session.End()
}

func (s *GithubSession) InitGithubAccessToken() {
	if *s.Options.GithubAccessToken == "" {
		accessToken := os.Getenv(AccessTokenEnvVariable)
		if accessToken == "" {
			s.Out.Fatal("No GitHub access token given. Please provide via command line option or in the %s environment variable.\n", AccessTokenEnvVariable)
		}
		s.GithubAccessToken = accessToken
	} else {
		s.GithubAccessToken = *s.Options.GithubAccessToken
	}
}

func (s *GithubSession) InitGithubClient() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	s.GithubClient = github.NewClient(tc)
	// s.GithubClient.UserAgent = fmt.Sprintf("%s v%s", Name, Version)
}

func (s *GithubSession) AddTarget(target *GithubOwner) {
	s.Lock()
	defer s.Unlock()
	for _, t := range s.Targets {
		if *target.ID == *t.ID {
			return
		}
	}
	s.Targets = append(s.Targets, target)
}

func (s *GithubSession) AddRepository(repository *GithubRepository) {
	s.Lock()
	defer s.Unlock()
	for _, r := range s.Repositories {
		if *repository.ID == *r.ID {
			return
		}
	}
	s.Repositories = append(s.Repositories, repository)
}

func NewGithubSession(options scan.Options) (*GithubSession, error) {
	var err error
	var githubRepos []*GithubRepository
	var targets []*GithubOwner
	sess := GithubSession{&session.Session{}, "", nil, targets, githubRepos}
	sess.Options = options
	err = session.ValidateNewSession(sess.Session)
	if err != nil {
		return nil, err
	}

	sess.Start()

	return &sess, nil
}
