package github

import (
	"context"
	"fmt"
	git2 "gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/git"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scan"
	"io/ioutil"
	"os"
	pathpkg "path"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
)

const (
	AccessTokenEnvVariable = "GITHUB_TOKEN"
	BASE                   = 10
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
	sess.Stats.Status = scan.StatusGathering
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
				sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), scan.Pluralize(len(repos), "repository", "repositories"), *target.Login)
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
	sess.Stats.Status = scan.StatusAnalyzing
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

	sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), scan.Pluralize(len(sess.Repositories), "repository", "repositories"))

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
				clone, dir, err := git2.CloneRepository(repo.CloneURL, repo.DefaultBranch, *sess.Options.CommitDepth)
				if err != nil {
					if err.Error() != "remote repository is empty" {
						sess.Out.Error("Error cloning repository %s: %s\n", *repo.FullName, err)
					}
					sess.Stats.IncrementRepositories()
					sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
					continue
				}
				sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", tid, *repo.FullName, dir)

				// // Path Scan
				// AnalyzeGithubRepoPaths(sess, repo, dir)
				// // Content Scan
				// AnalyzeGithubRepoContents(sess, repo, clone, dir)

				sess.Out.Debug("[THREAD #%d][%s] Fetching the checkpoint.\n", tid, *repo.FullName)
				checkpoint, err := scan.GetCheckpoint(strconv.FormatInt(*repo.ID, BASE), sess.Store.Connection)
				if err != nil {
					sess.Out.Debug("DB Error: %s\n", err)
				}

				if checkpoint == "" {
					//Scanning repo first time
					ScanGithubRepoCurrentRevision(sess, repo, dir)
				} else {
					ScanGithubRepoLatestCommits(sess, repo, clone, dir, checkpoint)
				}
				scan.UpdateCheckpoint(dir, strconv.FormatInt(*repo.ID, BASE), sess.Store.Connection)

				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.FullName)
				os.RemoveAll(dir)
				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.FullName, dir)
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

// ScanGithubRepoCurrentRevision runs the file scan for complete github repo.
// It scans only the lastest revision. rather than scanning the entire commit history
func ScanGithubRepoCurrentRevision(sess *GithubSession, repo *GithubRepository, dir string) {
	sess.Out.Debug(" Fetching %s repository files.\n", *repo.FullName)
	paths, err := scan.GatherPaths(dir, *repo.DefaultBranch)
	if err != nil {
		sess.Out.Error("Error while fetching the file paths of %s repository: %s\n", dir, err)
		return
	}
	sess.Out.Debug("[THREAD][%s] Fetching repository files of: %s\n", *repo.FullName, dir)
	for _, path := range paths {
		sess.Out.Debug("%s\n", path)
		content, err := ioutil.ReadFile(pathpkg.Join(dir, path))
		if err != nil {
			sess.Out.Error("[FILE NOT FOUND]: %s/%s\n", dir, path)
			continue
		}
		matchFile := scan.NewMatchFile(path, string(content))
		if matchFile.IsSkippable() {
			sess.Out.Debug("[THREAD][%s] Skipping %s\n", *repo.FullName, matchFile.Path)
			continue
		}
		sess.Out.Debug("[THREAD][%s] Matching: %s...\n", *repo.FullName, matchFile.Path)
		for _, signature := range scan.Signatures {
			if signature.Match(matchFile) {

				finding := &scan.Finding{
					FilePath:        path,
					Action:          signature.Part(),
					Description:     signature.Description(),
					Comment:         signature.Comment(),
					RepositoryOwner: *repo.Owner,
					RepositoryName:  *repo.Name,
					RepositoryUrl:   *repo.URL,
					FileUrl:         fmt.Sprintf("%s/blob/%s/%s", *repo.URL, *repo.DefaultBranch, matchFile.Path),
				}
				finding.Initialize()
				sess.AddFinding(finding)

				sess.Out.Warn(" %s: %s\n", strings.ToUpper(scan.PathScan), finding.Description)
				sess.Out.Info("  Repo.......: %s\n", *repo.FullName)
				sess.Out.Info("  Path.......: %s\n", finding.FilePath)
				sess.Out.Info("  Comment....: %s\n", finding.Comment)
				sess.Out.Info("  File URL...: %s\n", finding.FileUrl)
				sess.Out.Info(" ------------------------------------------------\n\n")
				sess.Stats.IncrementFindings()
			}
		}
	}
}

// ScanGithubRepoLatestCommits run a scan to analyze the diffs present in the commit history
// It will scan the commit history till the checkpoint (last scanned commit) is reached
func ScanGithubRepoLatestCommits(sess *GithubSession, repo *GithubRepository, clone *git.Repository, dir, checkpoint string) {
	history, err := git2.GetRepositoryHistory(clone)
	if err != nil {
		sess.Out.Error("[THREAD][%s] Error getting commit history: %s\n", *repo.FullName, err)
		return
	}
	sess.Out.Debug("[THREAD][%s] Number of commits: %d\n", *repo.FullName, len(history))

	for _, commit := range history {
		if strings.TrimSpace(commit.Hash.String()) == strings.TrimSpace(checkpoint) {
			//checkpoint reached
			break
		}
		sess.Out.Debug("[THREAD][%s] Analyzing commit: %s\n", *repo.FullName, commit.Hash)
		changes, _ := git2.GetChanges(commit, clone)
		sess.Out.Debug("[THREAD][%s] Changes in %s: %d\n", *repo.FullName, commit.Hash, len(changes))
		for _, change := range changes {
			path := git2.GetChangePath(change)
			allContent := ""
			sess.Out.Debug("FILE: %s/%s\n", dir, path)
			sess.Out.Debug("Repo URL: %s/commit/%s\n", *repo.URL, commit.Hash.String())
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
				sess.Out.Debug("[THREAD][%s] Skipping %s\n", *repo.FullName, matchFile.Path)
				continue
			}
			sess.Out.Debug("[THREAD][%s] Matching: %s...\n", *repo.FullName, matchFile.Path)
			for _, signature := range scan.Signatures {
				if signature.Match(matchFile) {
					// check if the matched signature is still present in the latest revision
					latestContent, err := ioutil.ReadFile(pathpkg.Join(dir, path))
					if err != nil {
						sess.Out.Info("[LATEST FILE NOT FOUND]: %s/%s\n", dir, path)
						continue
					}
					matchFile = scan.NewMatchFile(path, string(latestContent))
					if signature.Match(matchFile) {
						finding := &scan.Finding{
							FilePath:        path,
							Action:          scan.ContentScan,
							Description:     signature.Description(),
							Comment:         signature.Comment(),
							RepositoryOwner: *repo.Owner,
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

						sess.Out.Warn(" %s: %s\n", strings.ToUpper(scan.ContentScan), finding.Description)
						sess.Out.Info("  Path.......: %s\n", finding.FilePath)
						sess.Out.Info("  Repo.......: %s\n", *repo.FullName)
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
		sess.Out.Debug("[THREAD][%s] Done analyzing changes in %s\n", *repo.FullName, commit.Hash)
	}
}
