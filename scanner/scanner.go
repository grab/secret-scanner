/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package scanner

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/grab/secret-scanner/scanner/state"

	"github.com/grab/secret-scanner/scanner/findings"

	gitHandler "github.com/grab/secret-scanner/common/git"
	"github.com/grab/secret-scanner/scanner/gitprovider"
	"github.com/grab/secret-scanner/scanner/session"
	"github.com/grab/secret-scanner/scanner/signatures"
	"gopkg.in/src-d/go-git.v4"
)

// NewlineRegex ...
var NewlineRegex = regexp.MustCompile(`\r?\n`)

// Scan starts the scanning process
func Scan(sess *session.Session, gitProvider gitprovider.GitProvider) {
	if *sess.Options.LocalPath != "" {
		LocalGitScan(sess, gitProvider)
		sess.End()
		return
	}

	gatherRepositories(sess, gitProvider)

	sess.Stats.Status = session.StatusAnalyzing
	var ch = make(chan *gitprovider.Repository, len(sess.Repositories))
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

	var authMethod transport.AuthMethod
	// for github and gitlab, only personal access token is required, username can be a placeholder
	switch *sess.Options.GitProvider {
	case gitprovider.GithubName:
		authMethod = &http.BasicAuth{
			Username: "secretscanner",
			Password: *sess.Options.Token,
		}
	case gitprovider.GitlabName:
		authMethod = &http.BasicAuth{
			Username: "secretscanner",
			Password: *sess.Options.Token,
		}
	case gitprovider.BitbucketName:
		authMethod = &http.BasicAuth{
			Username: gitProvider.GetAdditionalParam(gitprovider.BitbucketParamUsername),
			Password: gitProvider.GetAdditionalParam(gitprovider.BitbucketParamPassword),
		}
	}

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

				// Clone repo
				sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", tid, repo.FullName)
				clone, cloneDir, err := gitHandler.CloneRepository(&repo.CloneURL, &repo.DefaultBranch, *sess.Options.CommitDepth, authMethod)
				if err != nil {
					if err.Error() != "Remote repository is empty" {
						sess.Out.Error("Error cloning repository %s: %s\n", repo.FullName, err)
					}
					sess.Stats.IncrementRepositories()
					sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
					continue
				}
				sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", tid, repo.FullName, cloneDir)

				// Get checkpoint
				sess.Out.Debug("[THREAD #%d][%s] Fetching the checkpoint.\n", tid, repo.FullName)
				checkpoint := ""

				if *sess.Options.State {
					latestHistory := sess.StateStore.Get(*sess.Options.GitProvider, repo.ID)
					if latestHistory != nil {
						checkpoint = latestHistory.CommitHash
					}
				}

				// Gather scan targets
				targets := sess.Options.ParseScanTargets()
				targetPaths, err := gitHandler.GatherPaths(cloneDir, repo.DefaultBranch, targets)
				if err != nil {
					sess.Out.Error("Failed to gather target paths for repo: %v", repo.FullName)
					return
				}

				targetPathMap := map[string]string{}
				for _, tp := range targetPaths {
					targetPathMap[path.Join(cloneDir, tp)] = tp
				}

				// Scan
				scanRevisions(sess, repo, clone, checkpoint, cloneDir, targetPathMap)
				latestCommitHash, err := gitHandler.GetLatestCommitHash(cloneDir)
				if err != nil {
					sess.Out.Error("Failed to get latest commit hash")
					return
				}

				if *sess.Options.State {
					err = sess.StateStore.Save(state.Create(*sess.Options.GitProvider, repo.ID, latestCommitHash, time.Now().String()))
					if err != nil {
						sess.Out.Error("Failed to save scan history: %v", err)
					}
				}

				// Cleanup
				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, repo.FullName)
				_ = os.RemoveAll(cloneDir)
				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, repo.FullName, cloneDir)
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

	sess.End()
}

// LocalGitScan starts a scan on local directory without first cloning from git provider
func LocalGitScan(sess *session.Session, gitProvider gitprovider.GitProvider) {
	sess.Stats.Status = session.StatusAnalyzing

	// Gather scan targets
	targets := sess.Options.ParseScanTargets()
	targetPaths, err := gitHandler.GatherPaths(*sess.Options.LocalPath, "master", targets)
	if err != nil {
		sess.Out.Error("Failed to gather target paths for repo: %v", *sess.Options.LocalPath)
		return
	}

	targetPathMap := map[string]string{}
	for _, tp := range targetPaths {
		stat, err := os.Stat(path.Join(*sess.Options.LocalPath, tp))
		if err != nil {
			continue
		}
		if stat.IsDir() {
			continue
		}
		targetPathMap[path.Join(*sess.Options.LocalPath, tp)] = tp
	}

	localID := fmt.Sprintf("%s/%s", strings.Trim(*sess.Options.LocalPath, "/"), strings.Trim(*sess.Options.ScanTarget, "/"))

	repo := &gitprovider.Repository{
		Owner:         "",
		ID:            localID,
		Name:          "",
		FullName:      *sess.Options.LocalPath,
		CloneURL:      "",
		URL:           "",
		DefaultBranch: "",
		Description:   "",
		Homepage:      "",
	}

	gitRepo, err := git.PlainOpen(*sess.Options.LocalPath)
	if err != nil {
		sess.Out.Error("Failed to open directory as git repo: %v", *sess.Options.LocalPath)
		return
	}

	// Get checkpoint
	checkpoint := ""
	if *sess.Options.State {
		latestHistory := sess.StateStore.Get(*sess.Options.GitProvider, repo.ID)
		if latestHistory != nil {
			checkpoint = latestHistory.CommitHash
		}
	}

	// Scan
	scanRevisions(sess, repo, gitRepo, checkpoint, *sess.Options.LocalPath, targetPathMap)

	sess.Stats.IncrementRepositories()
	sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))

	latestCommitHash, err := gitHandler.GetLatestCommitHash(*sess.Options.LocalPath)
	if err != nil {
		fmt.Println(err)
	}

	if *sess.Options.State {
		err = sess.StateStore.Save(state.Create(*sess.Options.GitProvider, repo.ID, latestCommitHash, time.Now().String()))
		if err != nil {
			sess.Out.Error("Failed to save scan history: %v", err)
		}
	}

	// NO cleanup for local scan
}

func gatherRepositories(sess *session.Session, gitProvider gitprovider.GitProvider) {
	var repos []*gitprovider.Repository

	if *sess.Options.Repos != "" {
		ids := strings.Split(*sess.Options.Repos, ",")
		for _, id := range ids {
			opt := map[string]string{}
			if gitProvider.Name() == gitprovider.GithubName || gitProvider.Name() == gitprovider.BitbucketName {
				idParts := strings.Split(id, "/")
				if len(idParts) != 2 {
					sess.Out.Error("Wrong Github option format (owner/repo): %v\n", errors.New("wrong Github option format"))
					continue
				}
				opt["owner"] = idParts[0]
				opt["repo"] = idParts[1]
			} else {
				opt["id"] = id
			}
			r, err := gitProvider.GetRepository(opt)
			if err != nil {
				sess.Out.Error("Error fetching the repo with ID %s: %s\n", id, err)
				continue
			}
			repos = append(repos, r)
		}
	}
	for _, repo := range repos {
		sess.Out.Info(" Retrieved repository: %s\n", repo.FullName)
		sess.AddRepository(repo)
	}
	sess.Stats.IncrementTargets()
	sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), Pluralize(len(repos), "repository", "repositories"), *sess.Options.GitProvider)
}

func scanRevisions(sess *session.Session, repo *gitprovider.Repository, clone *git.Repository, checkpoint, cloneDir string, targetPathMap map[string]string) {
	if checkpoint != "" {
		scanGitCommits(sess, repo, clone, cloneDir, checkpoint, targetPathMap)
	} else {
		scanCurrentGitRevision(sess, repo, cloneDir, targetPathMap)
	}
}

// scanCurrentGitRevision runs the file scan for complete gitlab repo.
// It scans only the lastest revision. rather than scanning the entire commit history
func scanCurrentGitRevision(sess *session.Session, repo *gitprovider.Repository, dir string, targetPathMap map[string]string) {
	sess.Out.Debug("[THREAD][%s] Fetching repository files of: %s\n", repo.FullName, dir)
	for absPath, subPath := range targetPathMap {
		sess.Out.Debug("Path: %s\n", absPath)
		content, err := ioutil.ReadFile(absPath)
		if err != nil {
			sess.Out.Error("[FILE NOT FOUND]: %s\n", absPath)
			continue
		}
		matchFile := signatures.NewMatchFile(subPath, string(content))
		if matchFile.IsSkippable() {
			sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
			continue
		}
		isTestContext := matchFile.IsTestContext()
		if isTestContext && *sess.Options.SkipTestContexts {
			sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
			continue
		}
		sess.Out.Debug("[THREAD][%s] Matching: %s...\n", repo.FullName, matchFile.Path)
		for _, signature := range sess.Signatures {
			matches := signature.Match(matchFile)

			for _, match := range matches {
				finding := &findings.Finding{
					FilePath:       subPath,
					Action:         signature.Part(),
					Description:    signature.Description(),
					Comment:        signature.Comment(),
					RepositoryName: repo.Name,
					RepositoryURL:  repo.URL,
					FileURL:        fmt.Sprintf("%s/blob/%s/%s", repo.URL, repo.DefaultBranch, subPath),
					Line:           match.Line,
					IsTestContext:  isTestContext,
				}

				if *sess.Options.LogSecret {
					finding.LineContent = match.LineContent
					finding.TruncateLineContent(findings.MaxLineChar)
				}

				hashID, err := finding.GenerateHashID()
				if err != nil {
					sess.Out.Error("Unable to generate hash ID for %v, skipping...", finding.FileURL)
					continue
				}
				finding.ID = hashID

				sess.AddFinding(finding)

				sess.Out.Warn(" %s: %s\n", strings.ToUpper(session.PathScan), finding.Description)
				sess.Out.Info("  Path........: %s\n", finding.FilePath)
				sess.Out.Info("  Repo........: %s\n", repo.FullName)
				sess.Out.Info("  Author......: %s\n", finding.CommitAuthor)
				sess.Out.Info("  Comment.....: %s\n", finding.Comment)
				sess.Out.Info("  File URL....: %s\n", finding.FileURL)
				sess.Out.Info("  Line........: %v\n", finding.Line)
				sess.Out.Info(" ------------------------------------------------\n\n")
				sess.Stats.IncrementFindings()
			}
		}
	}
}

// scanGitCommits run a scan to analyze the diffs present in the commit history
// It will scan the commit history till the checkpoint (last scanned commit) is reached
func scanGitCommits(sess *session.Session, repo *gitprovider.Repository, clone *git.Repository, dir, checkpoint string, targetPathMap map[string]string) {
	commitHistories, err := gitHandler.GetRepositoryHistory(clone)
	if err != nil {
		sess.Out.Error("[THREAD][%s] Error getting commit history: %s\n", repo.FullName, err)
		return
	}
	sess.Out.Debug("[THREAD][%s] Number of commits: %d\n", repo.FullName, len(commitHistories))

	for _, commit := range commitHistories {
		if strings.TrimSpace(commit.Hash.String()) == strings.TrimSpace(checkpoint) {
			sess.Out.Debug("\nCheckpoint Reached !!\n")
			break
		}
		sess.Out.Debug("[THREAD][%s] Analyzing commit: %s\n", repo.FullName, commit.Hash)
		changes, _ := gitHandler.GetChanges(commit, clone)
		sess.Out.Debug("[THREAD][%s] Changes in %s: %d\n", repo.FullName, commit.Hash, len(changes))
		for _, change := range changes {
			p := gitHandler.GetChangePath(change)

			_, exists := targetPathMap[path.Join(dir, p)]
			if len(targetPathMap) > 0 && !exists {
				continue
			}

			allContent := ""
			sess.Out.Debug("FILE: %s/%s\n", dir, p)
			sess.Out.Debug("Repo URL: %s/commit/%s\n", repo.URL, commit.Hash.String())
			patch, _ := gitHandler.GetPatch(change)
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
			matchFile := signatures.NewMatchFile(p, allContent)
			if matchFile.IsSkippable() {
				sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
				continue
			}
			isTestContext := matchFile.IsTestContext()
			if isTestContext && *sess.Options.SkipTestContexts {
				sess.Out.Debug("[THREAD][%s] Skipping %s\n", repo.FullName, matchFile.Path)
				continue
			}
			sess.Out.Debug("[THREAD][%s] Matching: %s...\n", repo.FullName, matchFile.Path)
			for _, signature := range sess.Signatures {
				matches := signature.Match(matchFile)

				for _, match := range matches {
					finding := &findings.Finding{
						FilePath:       p,
						Action:         signature.Part(),
						Description:    signature.Description(),
						Comment:        signature.Comment(),
						RepositoryName: repo.Name,
						RepositoryURL:  repo.URL,
						FileURL:        fmt.Sprintf("%s/blob/%s/%s", repo.URL, repo.DefaultBranch, p),
						Line:           match.Line,
						IsTestContext:  isTestContext,
					}

					if *sess.Options.LogSecret {
						finding.LineContent = match.LineContent
						finding.TruncateLineContent(findings.MaxLineChar)
					}

					hashID, err := finding.GenerateHashID()
					if err != nil {
						sess.Out.Error("Unable to generate hash ID for %v, skipping...", finding.FileURL)
						continue
					}
					finding.ID = hashID

					sess.AddFinding(finding)

					sess.Out.Warn(" %s: %s\n", strings.ToUpper(session.PathScan), finding.Description)
					sess.Out.Info("  Path........: %s\n", finding.FilePath)
					sess.Out.Info("  Repo........: %s\n", repo.FullName)
					sess.Out.Info("  Author......: %s\n", finding.CommitAuthor)
					sess.Out.Info("  Comment.....: %s\n", finding.Comment)
					sess.Out.Info("  File URL....: %s\n", finding.FileURL)
					sess.Out.Info("  Line........: %v\n", finding.Line)
					sess.Out.Info(" ------------------------------------------------\n\n")
					sess.Stats.IncrementFindings()
				}

				//if signature.Match(matchFile) {
				//	latestContent, err := ioutil.ReadFile(path.Join(dir, p))
				//	if err != nil {
				//		sess.Out.Info("[LATEST FILE NOT FOUND]: %s/%s\n", dir, p)
				//		continue
				//	}
				//	matchFile = signatures.NewMatchFile(p, string(latestContent))
				//	if signature.Match(matchFile) {
				//		finding := &findings.Finding{
				//			FilePath:       p,
				//			Action:         session.ContentScan,
				//			Description:    signature.Description(),
				//			Comment:        signature.Comment(),
				//			RepositoryName: repo.Name,
				//			CommitHash:     commit.Hash.String(),
				//			CommitMessage:  strings.TrimSpace(commit.Message),
				//			CommitAuthor:   commit.Author.String(),
				//			RepositoryURL:  repo.URL,
				//			FileURL:        fmt.Sprintf("%s/blob/%s/%s", repo.URL, repo.DefaultBranch, p),
				//			CommitURL:      fmt.Sprintf("%s/commit/%s", repo.URL, commit.Hash.String()),
				//			IsTestContext:  isTestContext,
				//		}
				//
				//		hashID, err := finding.GenerateHashID()
				//		if err != nil {
				//			sess.Out.Error("Unable to generate hash ID for %v, skipping...", finding.FileURL)
				//			continue
				//		}
				//		finding.ID = hashID
				//
				//		sess.AddFinding(finding)
				//
				//		sess.Out.Warn(" %s: %s\n", strings.ToUpper(session.ContentScan), finding.Description)
				//		sess.Out.Info("  Path.......: %s\n", finding.FilePath)
				//		sess.Out.Info("  Repo.......: %s\n", repo.FullName)
				//		sess.Out.Info("  Message....: %s\n", TruncateString(finding.CommitMessage, 100))
				//		sess.Out.Info("  Author.....: %s\n", finding.CommitAuthor)
				//		sess.Out.Info("  Comment....: %s\n", finding.Comment)
				//		sess.Out.Info("  File URL...: %s\n", finding.FileURL)
				//		sess.Out.Info("  Commit URL.: %s\n", finding.CommitURL)
				//		sess.Out.Info(" ------------------------------------------------\n\n")
				//		sess.Stats.IncrementFindings()
				//	}
				//}
			}
			sess.Stats.IncrementFiles()
		}
		sess.Stats.IncrementCommits()
		sess.Out.Debug("[THREAD][%s] Done analyzing changes in %s\n", repo.FullName, commit.Hash)
	}
}

// Pluralize makes word plural
func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// TruncateString truncates string from whitespace pre/suf-fix
func TruncateString(str string, maxLength int) string {
	str = NewlineRegex.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	if len(str) > maxLength {
		str = fmt.Sprintf("%s...", str[0:maxLength])
	}
	return str
}
