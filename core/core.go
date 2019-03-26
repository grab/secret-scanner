package core

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func GatherPaths(dir, branch string) ([]string, error) {
	os.Chdir(dir)
	gitcmd := "git"
	listTree := "ls-tree"
	op1 := "-r"
	op2 := "--name-only"
	out, err := exec.Command(gitcmd, listTree, op1, branch, op2).CombinedOutput()
	if err != nil {
		return nil, err
	}
	cmdout := fmt.Sprintf("%s", string(out))
	paths := strings.Split(cmdout, "\n")
	return paths, nil
}
