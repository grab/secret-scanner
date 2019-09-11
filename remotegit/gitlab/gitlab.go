package gitlab

import (
	"fmt"
	git2 "gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/git"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/logic/scan"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/remotegit"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
	"io/ioutil"
	"os"
	pathpkg "path"
	"strconv"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
)

const (
	GitlabTokenEnvVariable = "GITLAB_TOKEN"
	GitlabEndpoint         = "https://gitlab.myteksi.net"
)

func GetAllRepositories(git *gitlab.Client) ([]*remotegit.Repository, error) {
	var allRepos []*remotegit.Repository
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}
	for {
		projects, resp, err := git.Projects.ListProjects(opt)
		if err != nil {
			return allRepos, err
		}
		for _, proj := range projects {
			r := remotegit.Repository{
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
			allRepos = append(allRepos, &r)
		}
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func GetRepository(git *gitlab.Client, id string) (*remotegit.Repository, error) {
	proj, _, err := git.Projects.GetProject(id, nil)
	if err != nil {
		return nil, err
	}
	repo := &remotegit.Repository{
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

func GatherGitlabRepos(sess *GitlabSession) {
	var repos []*remotegit.Repository
	var err error
	if *sess.Options.Repos != "" {
		//Fetching the repos prodided in repo-list
		if !scan.FileExists(*sess.Options.Repos) {
			sess.Out.Error(" No such file exists in: %s\n", *sess.Options.Repos)
		}
		data, err := ioutil.ReadFile(*sess.Options.Repos)
		if err != nil {
			sess.Out.Error(" Failed to load the repo list provided: %s\n", err)
		}
		ids := strings.Split(string(data), ",")
		for _, id := range ids {
			r, err := GetRepository(sess.GitlabClient, id)
			if err != nil {
				sess.Out.Error("Error fetching the repo with ID %s: %s\n", id, err)
				continue
			}
			repos = append(repos, r)
		}
	} else {
		//fetch all repos
		repos, err = GetAllRepositories(sess.GitlabClient)
		if err != nil {
			sess.Out.Error(" Failed to retrieve repositories: %s\n", err)
		}
	}
	for _, repo := range repos {
		sess.Out.Info(" Retrieved repository: %s\n", repo.FullName)
		sess.AddGitlabRepository(repo)
	}
	sess.Stats.IncrementTargets()
	sess.Out.Info(" Retrieved %d %s from GITLAB\n", len(repos), scan.Pluralize(len(repos), "repository", "repositories"))
}

func AnalyzeGitlabRepositories(sess *GitlabSession) {
	sess.Stats.Status = session.StatusAnalyzing
	var ch = make(chan *remotegit.Repository, len(sess.GitlabRepos))
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

	sess.Out.Important("Analyzing %d %s...\n", len(sess.GitlabRepos), scan.Pluralize(len(sess.GitlabRepos), "repository", "repositories"))

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

				sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", tid, repo.FullName)
				clone, dir, err := git2.CloneRepository(&repo.CloneURL, &repo.DefaultBranch, *sess.Options.CommitDepth)
				if err != nil {
					if err.Error() != "Remote repository is empty" {
						sess.Out.Error("Error cloning repository %s: %s\n", repo.FullName, err)
					}
					sess.Stats.IncrementRepositories()
					sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.GitlabRepos))
					continue
				}
				sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", tid, repo.FullName, dir)

				sess.Out.Debug("[THREAD #%d][%s] Fetching the checkpoint.\n", tid, repo.FullName)
				checkpoint, err := scan.GetCheckpoint(strconv.Itoa(int(repo.ID)), sess.Store.Connection)
				if err != nil {
					sess.Out.Debug("DB Error: %s\n", err)
				}

				if checkpoint == "" {
					//Scanning repo first time
					ScanGitlabRepoCurrentRevision(sess, repo, dir)
				} else {
					ScanGitlabRepoLatestCommits(sess, repo, clone, dir, checkpoint)
				}
				scan.UpdateCheckpoint(dir, strconv.Itoa(int(repo.ID)), sess.Store.Connection)

				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, repo.FullName)
				os.RemoveAll(dir)
				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, repo.FullName, dir)
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

// ScanGitlabRepoCurrentRevision runs the file scan for complete gitlab repo.
// It scans only the lastest revision. rather than scanning the entire commit history
func ScanGitlabRepoCurrentRevision(sess *GitlabSession, repo *remotegit.Repository, dir string) {
	paths, err := scan.GatherPaths(dir, repo.DefaultBranch)
	if err != nil {
		sess.Out.Error("Error while fetching the file paths of %s repository: %s\n", dir, err)
		return
	}
	sess.Out.Debug("[THREAD][%s] Fetching repository files of: %s\n", repo.FullName, dir)
	for _, path := range paths {
		sess.Out.Debug("Path: %s\n", path)
		content, err := ioutil.ReadFile(pathpkg.Join(dir, path))
		if err != nil {
			sess.Out.Error("[FILE NOT FOUND]: %s/%s\n", dir, path)
			continue
		}
		matchFile := scan.NewMatchFile(path, string(content))
		if matchFile.IsSkippable() {
			sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
			continue
		}
		sess.Out.Debug("[THREAD][%s] Matching: %s...\n", repo.FullName, matchFile.Path)
		for _, signature := range scan.Signatures {
			if signature.Match(matchFile) {
				finding := &scan.Finding{
					FilePath:       path,
					Action:         signature.Part(),
					Description:    signature.Description(),
					Comment:        signature.Comment(),
					RepositoryName: repo.Name,
					RepositoryUrl:  repo.URL,
					FileUrl:        fmt.Sprintf("%s/blob/%s/%s", repo.URL, repo.DefaultBranch, path),
				}
				finding.Initialize()
				sess.AddFinding(finding)

				sess.Out.Warn(" %s: %s\n", strings.ToUpper(session.PathScan), finding.Description)
				sess.Out.Info("  Path.......: %s\n", finding.FilePath)
				sess.Out.Info("  Repo.......: %s\n", repo.FullName)
				sess.Out.Info("  Author.....: %s\n", finding.CommitAuthor)
				sess.Out.Info("  Comment....: %s\n", finding.Comment)
				sess.Out.Info("  File URL...: %s\n", finding.FileUrl)
				sess.Out.Info(" ------------------------------------------------\n\n")
				sess.Stats.IncrementFindings()
			}
		}
	}
}

// ScanGitlabRepoLatestCommits run a scan to analyze the diffs present in the commit history
// It will scan the commit history till the checkpoint (last scanned commit) is reached
func ScanGitlabRepoLatestCommits(sess *GitlabSession, repo *remotegit.Repository, clone *git.Repository, dir, checkpoint string) {
	history, err := git2.GetRepositoryHistory(clone)
	if err != nil {
		sess.Out.Error("[THREAD][%s] Error getting commit history: %s\n", repo.FullName, err)
		return
	}
	sess.Out.Debug("[THREAD][%s] Number of commits: %d\n", repo.FullName, len(history))

	for _, commit := range history {
		if strings.TrimSpace(commit.Hash.String()) == strings.TrimSpace(checkpoint) {
			sess.Out.Debug("\nCheckpoint Reached !!\n")
			break
		}
		sess.Out.Debug("[THREAD][%s] Analyzing commit: %s\n", repo.FullName, commit.Hash)
		changes, _ := git2.GetChanges(commit, clone)
		sess.Out.Debug("[THREAD][%s] Changes in %s: %d\n", repo.FullName, commit.Hash, len(changes))
		for _, change := range changes {
			path := git2.GetChangePath(change)
			allContent := ""
			sess.Out.Debug("FILE: %s/%s\n", dir, path)
			sess.Out.Debug("Repo URL: %s/commit/%s\n", repo.URL, commit.Hash.String())
			patch, _ := git2.GetPatch(change)
			diffs := patch.FilePatches()
			for _, diff := range diffs {
				chunks := diff.Chunks()
				for _, chunk := range chunks {
					if chunk.Type() == 1 {
						allContent += chunk.Content()
						allContent += "\n\n"
					}
				}
			}
			matchFile := scan.NewMatchFile(path, allContent)
			if matchFile.IsSkippable() {
				sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
				continue
			}
			sess.Out.Debug("[THREAD][%s] Matching: %s...\n", repo.FullName, matchFile.Path)
			for _, signature := range scan.Signatures {
				if signature.Match(matchFile) {
					latestContent, err := ioutil.ReadFile(pathpkg.Join(dir, path))
					if err != nil {
						sess.Out.Info("[LATEST FILE NOT FOUND]: %s/%s\n", dir, path)
						continue
					}
					matchFile = scan.NewMatchFile(path, string(latestContent))
					if signature.Match(matchFile) {
						finding := &scan.Finding{
							FilePath:       path,
							Action:         session.ContentScan,
							Description:    signature.Description(),
							Comment:        signature.Comment(),
							RepositoryName: repo.Name,
							CommitHash:     commit.Hash.String(),
							CommitMessage:  strings.TrimSpace(commit.Message),
							CommitAuthor:   commit.Author.String(),
							RepositoryUrl:  repo.URL,
							FileUrl:        fmt.Sprintf("%s/blob/%s/%s", repo.URL, repo.DefaultBranch, path),
							CommitUrl:      fmt.Sprintf("%s/commit/%s", repo.URL, commit.Hash.String()),
						}
						finding.Initialize()
						sess.AddFinding(finding)

						sess.Out.Warn(" %s: %s\n", strings.ToUpper(session.ContentScan), finding.Description)
						sess.Out.Info("  Path.......: %s\n", finding.FilePath)
						sess.Out.Info("  Repo.......: %s\n", repo.FullName)
						sess.Out.Info("  Message....: %s\n", scan.TruncateString(finding.CommitMessage, 100))
						sess.Out.Info("  Author.....: %s\n", finding.CommitAuthor)
						sess.Out.Info("  Comment....: %s\n", finding.Comment)
						sess.Out.Info("  File URL...: %s\n", finding.FileUrl)
						sess.Out.Info("  Commit URL.: %s\n", finding.CommitUrl)
						sess.Out.Info(" ------------------------------------------------\n\n")
						sess.Stats.IncrementFindings()
					}
				}
			}
			sess.Stats.IncrementFiles()
		}
		sess.Stats.IncrementCommits()
		sess.Out.Debug("[THREAD][%s] Done analyzing changes in %s\n", repo.FullName, commit.Hash)
	}
}
