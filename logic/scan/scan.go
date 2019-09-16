package scan

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var NewlineRegex = regexp.MustCompile(`\r?\n`)

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func TruncateString(str string, maxLength int) string {
	str = NewlineRegex.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	if len(str) > maxLength {
		str = fmt.Sprintf("%s...", str[0:maxLength])
	}
	return str
}

func GatherPaths(dir, branch string, targets []string) ([]string, error) {
	os.Chdir(dir)
	gitcmd := "git"
	listTree := "ls-tree"
	op1 := "-r"
	op2 := "--name-only"
	var paths []string

	if len(targets) == 0 {
		out, err := exec.Command(gitcmd, listTree, op1, branch, op2).CombinedOutput()
		if err != nil {
			return nil, err
		}
		cmdout := fmt.Sprintf("%s", string(out))
		paths = append(paths, strings.Split(cmdout, "\n")...)
	}

	for _, t := range targets {
		out, err := exec.Command(gitcmd, listTree, op1, fmt.Sprintf("%s:%s", branch, t), op2).CombinedOutput()
		if err != nil {
			return nil, err
		}
		cmdout := fmt.Sprintf("%s", string(out))
		currentPaths := strings.Split(cmdout, "\n")
		for i, p := range currentPaths {
			currentPaths[i] = path.Join(t, p)
		}
		paths = append(paths, currentPaths...)
	}
	return paths, nil
}

// GetLatestCommitHash runs a git cmd to return latest commit hash
func GetLatestCommitHash(dir string) (string, error) {
	os.Chdir(dir)
	gitcmd := "git"
	task := "rev-parse"
	op1 := "--verify"
	op2 := "HEAD"
	out, err := exec.Command(gitcmd, task, op1, op2).CombinedOutput()
	if err != nil {
		return "", err
	}
	commitHash := fmt.Sprintf("%s", string(out))
	return commitHash, nil
}

// GetCheckpoint returns the last scanned commit hash
func GetCheckpoint(repoId string, db *sql.DB) (string, error) {
	var commitHash sql.NullString
	query := sq.Select("commit_hash").
		From("scan_history").
		Where("repo_id = ?", repoId).
		RunWith(db).
		QueryRow()

	err := query.Scan(&commitHash)
	if err != nil {
		return "", err
	}

	return commitHash.String, nil
}

// UpdateCheckpoint insert (or updates if exists) the repo id
// and its latest commit hash in the DB
func UpdateCheckpoint(dir, repoId string, db *sql.DB) error {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	latestCommitHash, err := GetLatestCommitHash(dir)
	if err != nil {
		return err
	}
	_, err = sq.Insert("scan_history").
		Columns("repo_id", "commit_hash").
		Values(repoId, latestCommitHash).
		Suffix("ON DUPLICATE KEY UPDATE commit_hash = ?", latestCommitHash).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
