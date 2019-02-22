package core

import (
  "os"
  "sync"
  "strings"
  "context"
  "github.com/google/go-github/github"
)

const (
  AccessTokenEnvVariable = "GITHUB_TOKEN"
)

type GithubOwner struct {
  Login     *string
  ID        *int64
  Type      *string
  Name      *string
  AvatarURL *string
  URL       *string
  Company   *string
  Blog      *string
  Location  *string
  Email     *string
  Bio       *string
}

type GithubRepository struct {
  Owner         *string
  ID            *int64
  Name          *string
  FullName      *string
  CloneURL      *string
  URL           *string
  DefaultBranch *string
  Description   *string
  Homepage      *string
}

func GetUserOrOrganization(login string, client *github.Client) (*GithubOwner, error) {
  ctx := context.Background()
  user, _, err := client.Users.Get(ctx, login)
  if err != nil {
    return nil, err
  }
  return &GithubOwner{
    Login:     user.Login,
    ID:        user.ID,
    Type:      user.Type,
    Name:      user.Name,
    AvatarURL: user.AvatarURL,
    URL:       user.HTMLURL,
    Company:   user.Company,
    Blog:      user.Blog,
    Location:  user.Location,
    Email:     user.Email,
    Bio:       user.Bio,
  }, nil
}

func GetRepositoriesFromOwner(login *string, client *github.Client) ([]*GithubRepository, error) {
  var allRepos []*GithubRepository
  loginVal := *login
  ctx := context.Background()
  opt := &github.RepositoryListOptions{
    Type: "sources",
  }

  for {
    repos, resp, err := client.Repositories.List(ctx, loginVal, opt)
    if err != nil {
      return allRepos, err
    }
    for _, repo := range repos {
      if !*repo.Fork {
        r := GithubRepository{
          Owner:         repo.Owner.Login,
          ID:            repo.ID,
          Name:          repo.Name,
          FullName:      repo.FullName,
          CloneURL:      repo.CloneURL,
          URL:           repo.HTMLURL,
          DefaultBranch: repo.DefaultBranch,
          Description:   repo.Description,
          Homepage:      repo.Homepage,
        }
        allRepos = append(allRepos, &r)
      }
    }
    if resp.NextPage == 0 {
      break
    }
    opt.Page = resp.NextPage
  }

  return allRepos, nil
}

func GetOrganizationMembers(login *string, client *github.Client) ([]*GithubOwner, error) {
  var allMembers []*GithubOwner
  loginVal := *login
  ctx := context.Background()
  opt := &github.ListMembersOptions{}
  for {
    members, resp, err := client.Organizations.ListMembers(ctx, loginVal, opt)
    if err != nil {
      return allMembers, err
    }
    for _, member := range members {
      allMembers = append(allMembers, &GithubOwner{Login: member.Login, ID: member.ID, Type: member.Type})
    }
    if resp.NextPage == 0 {
      break
    }
    opt.Page = resp.NextPage
  }
  return allMembers, nil
}

func GatherTargets(sess *GithubSession) {
  sess.Stats.Status = StatusGathering
  sess.Out.Important("Gathering targets...\n")
  for _, login := range sess.Options.Logins {
    target, err := GetUserOrOrganization(login, sess.GithubClient)
    if err != nil {
      sess.Out.Error(" Error retrieving information on %s: %s\n", login, err)
      continue
    }
    sess.Out.Debug("%s (ID: %d) type: %s\n", *target.Login, *target.ID, *target.Type)
    sess.AddTarget(target)
    if *sess.Options.NoExpandOrgs == false && *target.Type == "Organization" {
      sess.Out.Debug("Gathering members of %s (ID: %d)...\n", *target.Login, *target.ID)
      members, err := GetOrganizationMembers(target.Login, sess.GithubClient)
      if err != nil {
        sess.Out.Error(" Error retrieving members of %s: %s\n", *target.Login, err)
        continue
      }
      for _, member := range members {
        sess.Out.Debug("Adding organization member %s (ID: %d) to targets\n", *member.Login, *member.ID)
        sess.AddTarget(member)
      }
    }
  }
}

func GatherRepositories(sess *GithubSession) {
  var ch = make(chan *GithubOwner, len(sess.Targets))
  var wg sync.WaitGroup
  var threadNum int
  if len(sess.Targets) == 1 {
    threadNum = 1
  } else if len(sess.Targets) <= *sess.Options.Threads {
    threadNum = len(sess.Targets) - 1
  } else {
    threadNum = *sess.Options.Threads
  }
  wg.Add(threadNum)
  sess.Out.Debug("Threads for repository gathering: %d\n", threadNum)
  for i := 0; i < threadNum; i++ {
    go func() {
      for {
        target, ok := <-ch
        if !ok {
          wg.Done()
          return
        }
        repos, err := GetRepositoriesFromOwner(target.Login, sess.GithubClient)
        if err != nil {
          sess.Out.Error(" Failed to retrieve repositories from %s: %s\n", *target.Login, err)
        }
        if len(repos) == 0 {
          continue
        }
        for _, repo := range repos {
          sess.Out.Debug(" Retrieved repository: %s\n", *repo.FullName)
          sess.AddRepository(repo)
        }
        sess.Stats.IncrementTargets()
        sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), Pluralize(len(repos), "repository", "repositories"), *target.Login)
      }
    }()
  }

  for _, target := range sess.Targets {
    ch <- target
  }
  close(ch)
  wg.Wait()
}

func AnalyzeRepositories(sess *GithubSession) {
	sess.Stats.Status = StatusAnalyzing
	var ch = make(chan *GithubRepository, len(sess.Repositories))
	var wg sync.WaitGroup
	var threadNum int
	if len(sess.Repositories) <= 1 {
		threadNum = 1
	} else if len(sess.Repositories) <= *sess.Options.Threads {
		threadNum = len(sess.Repositories) - 1
	} else {
		threadNum = *sess.Options.Threads
	}
	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository analysis: %d\n", threadNum)
	
	sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), Pluralize(len(sess.Repositories), "repository", "repositories"))
	
	for i := 0; i < threadNum; i++ {
		go func(tid int) {
		for {
			sess.Out.Debug("[THREAD #%d] Requesting new repository to analyze...\n", tid)
			repo, ok := <-ch
			if !ok {
			sess.Out.Debug("[THREAD #%d] No more tasks, marking WaitGroup as done\n", tid)
			wg.Done()
			return
			}
	
			sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", tid, *repo.FullName)
			clone, path, err := CloneRepository(repo.CloneURL, repo.DefaultBranch, *sess.Options.CommitDepth)
			if err != nil {
			if err.Error() != "remote repository is empty" {
				sess.Out.Error("Error cloning repository %s: %s\n", *repo.FullName, err)
			}
			sess.Stats.IncrementRepositories()
			sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
			continue
			}
			sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", tid, *repo.FullName, path)
	
			history, err := GetRepositoryHistory(clone)
			if err != nil {
			sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", tid, *repo.FullName, err)
			os.RemoveAll(path)
			sess.Stats.IncrementRepositories()
			sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
			continue
			}
			sess.Out.Debug("[THREAD #%d][%s] Number of commits: %d\n", tid, *repo.FullName, len(history))
	
			for _, commit := range history {
				sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.FullName, commit.Hash)
				changes, _ := GetChanges(commit, clone)
				sess.Out.Debug("[THREAD #%d][%s] Changes in %s: %d\n", tid, *repo.FullName, commit.Hash, len(changes))
				for _, change := range changes {
					changeAction := GetChangeAction(change)
					path := GetChangePath(change)
					matchFile := NewMatchFile(path)
					if matchFile.IsSkippable() {
					sess.Out.Debug("[THREAD #%d][%s] Skipping %s\n", tid, *repo.FullName, matchFile.Path)
					continue
					}
					sess.Out.Debug("[THREAD #%d][%s] Matching: %s...\n", tid, *repo.FullName, matchFile.Path)
					for _, signature := range Signatures {
					if signature.Match(matchFile) {
		
						finding := &Finding{
						FilePath:        path,
						Action:          changeAction,
						Description:     signature.Description(),
						Comment:         signature.Comment(),
						RepositoryOwner: *repo.Owner,
						RepositoryName:  *repo.Name,
						CommitHash:      commit.Hash.String(),
						CommitMessage:   strings.TrimSpace(commit.Message),
						CommitAuthor:    commit.Author.String(),
						}
						finding.Initialize()
						sess.AddFinding(finding)
		
						sess.Out.Warn(" %s: %s\n", strings.ToUpper(changeAction), finding.Description)
						sess.Out.Info("  Path.......: %s\n", finding.FilePath)
						sess.Out.Info("  Repo.......: %s\n", *repo.FullName)
						sess.Out.Info("  Message....: %s\n", TruncateString(finding.CommitMessage, 100))
						sess.Out.Info("  Author.....: %s\n", finding.CommitAuthor)
						if finding.Comment != "" {
						sess.Out.Info("  Comment....: %s\n", finding.Comment)
						}
						sess.Out.Info("  File URL...: %s\n", finding.FileUrl)
						sess.Out.Info("  Commit URL.: %s\n", finding.CommitUrl)
						sess.Out.Info(" ------------------------------------------------\n\n")
						sess.Stats.IncrementFindings()
						break
					}
					}
					sess.Stats.IncrementFiles()
				}
				sess.Stats.IncrementCommits()
				sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.FullName, commit.Hash)
			}
			sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.FullName)
			os.RemoveAll(path)
			sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.FullName, path)
			sess.Stats.IncrementRepositories()
			sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
		}
		}(i)
	}
	for _, repo := range sess.Repositories {
		ch <- repo
	}
	close(ch)
	wg.Wait()
}
