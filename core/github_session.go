package core

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubSession struct {
	*Session
	GithubAccessToken string         `json:"-"`
	GithubClient      *github.Client `json:"-"`
	Targets           []*GithubOwner
	Repositories      []*GithubRepository
}

func (s *GithubSession) Start() {
	s.Session.Start()
	s.InitGithubAccessToken()
	s.InitGithubClient()
	// s.InitRouter()
}

func (s *GithubSession) Finish() {
	s.Session.Finish()
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

func NewGithubSession(options Options) (*GithubSession, error) {
	var err error
	var githubRepos []*GithubRepository
	var targets []*GithubOwner
	session := GithubSession{&Session{}, "", nil, targets, githubRepos}
	session.Options = options
	err = ValidateNewSession(session.Session)
	if err != nil {
		return nil, err
	}

	session.Start()

	return &session, nil
}
