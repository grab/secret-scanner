package core

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

type GitlabRepository struct {
	Owner         *string
	ID            *int
	Name          *string
	FullName      *string
	CloneURL      *string
	URL           *string
	DefaultBranch *string
	Description   *string
	Homepage      *string
}

const (
	GitlabTokenEnvVariable = "GITLAB_TOKEN"
	GitlabEndpoint         = "https://gitlab.myteksi.net"
)

func GetAllRepositories(git *gitlab.Client) ([]*GitlabRepository, error) {
	var allRepos []*GitlabRepository
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}
	// projects, _, err := git.Projects.ListProjects(opt)
	// 	if err != nil {
	// 		return allRepos, err
	// 	}
	// 	for _, proj := range projects {
	// 		r := GitlabRepository{
	// 			ID:            &proj.ID,
	// 			Name:          &proj.Name,
	// 			FullName:      &proj.Name,
	// 			CloneURL:      &proj.SSHURLToRepo,
	// 			URL:           &proj.WebURL,
	// 			DefaultBranch: &proj.DefaultBranch,
	// 			Description:   &proj.Description,
	// 			Homepage:      &proj.WebURL,
	// 			Owner:         nil,
	// 			}
	// 		allRepos = append(allRepos, &r)
	// 		log.Printf("Found project: %s", proj.Name)
	// 	}
	for {
		projects, resp, err := git.Projects.ListProjects(opt)
		if err != nil {
			return allRepos, err
		}
		for _, proj := range projects {
			r := GitlabRepository{
				ID:            &proj.ID,
				Name:          &proj.Name,
				FullName:      &proj.Name,
				CloneURL:      &proj.SSHURLToRepo,
				URL:           &proj.WebURL,
				DefaultBranch: &proj.DefaultBranch,
				Description:   &proj.Description,
				Homepage:      &proj.WebURL,
				Owner:         nil,
			}
			allRepos = append(allRepos, &r)
		}
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func GatherGitlabRepos(sess *GitlabSession) {
	repos, err := GetAllRepositories(sess.GitlabClient)
	if err != nil {
		sess.Out.Error(" Failed to retrieve repositories: %s\n", err)
	}
	for _, repo := range repos {
		sess.Out.Debug(" Retrieved repository: %s\n", *repo.FullName)
		sess.AddGitlabRepository(repo)
	}
	sess.Stats.IncrementTargets()
	sess.Out.Info(" Retrieved %d %s from GITLAB\n", len(repos), Pluralize(len(repos), "repository", "repositories"))
}

func AnalyzeGitlabRepositories(sess *GitlabSession) {
	sess.Stats.Status = StatusAnalyzing
	var ch = make(chan *GitlabRepository, len(sess.GitlabRepos))
	var wg sync.WaitGroup
	var threadNum int
	if len(sess.GitlabRepos) <= 1 {
		threadNum = 1
	} else if len(sess.GitlabRepos) <= *sess.Options.Threads {
		threadNum = len(sess.GitlabRepos) - 1
	} else {
		threadNum = *sess.Options.Threads
	}
	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository analysis: %d\n", threadNum)

	sess.Out.Important("Analyzing %d %s...\n", len(sess.GitlabRepos), Pluralize(len(sess.GitlabRepos), "repository", "repositories"))

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
					sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.GitlabRepos))
					continue
				}
				sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", tid, *repo.FullName, path)

				history, err := GetRepositoryHistory(clone)
				if err != nil {
					sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", tid, *repo.FullName, err)
					os.RemoveAll(path)
					sess.Stats.IncrementRepositories()
					sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.GitlabRepos))
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
									RepositoryOwner: "",
									RepositoryName:  *repo.Name,
									CommitHash:      commit.Hash.String(),
									CommitMessage:   strings.TrimSpace(commit.Message),
									CommitAuthor:    commit.Author.String(),
									RepositoryUrl:   *repo.URL,
									FileUrl:         fmt.Sprintf("%s/blob/%s/%s", *repo.URL, commit.Hash.String(), path),
									CommitUrl:       fmt.Sprintf("%s/commit/%s", *repo.URL, commit.Hash.String()),
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
				sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.GitlabRepos))
			}
		}(i)
	}
	for _, repo := range sess.GitlabRepos {
		ch <- repo
	}
	close(ch)
	wg.Wait()
}
