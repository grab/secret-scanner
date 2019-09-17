package options

import (
	"flag"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"os"
	"strings"
)

type Options struct {
	CommitDepth        *int
	Threads            *int
	Save               *string `json:"-"`
	Load               *string `json:"-"`
	Silent             *bool
	Debug              *bool

	GitProvider        *string
	GitProviderBaseURL *string
	GitProviderToken   *string
	EnvFilePath        *string
	RepoID             *string
	ScanTargets        *string
	Repos              *string
	GitScanPath        *string
}

func (o Options) ValidateOptions() error {
	if *o.GitProvider != gitprovider.GithubName && *o.GitProvider != gitprovider.GitlabName {
		return ErrInvalidGitProvider
	}

	if *o.RepoID != "" && *o.Repos != "" {
		return ErrRepoOptionConflict
	}
	if *o.EnvFilePath != "" {
		if _, err := os.Stat(*o.EnvFilePath); os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (o *Options) Parse() []string {
	return strings.Split(*o.ScanTargets, ",")
}

func (o *Options) ParseScanTargets() []string {
	return strings.Split(*o.ScanTargets, ",")
}

func Parse() (Options, error) {
	options := Options{
		CommitDepth:        flag.Int("commit-depth", 500, "Number of repository commits to process"),
		Threads:            flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Save:               flag.String("save", "", "Save session to file"),
		Load:               flag.String("load", "", "Load session file"),
		Silent:             flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:              flag.Bool("debug", false, "Print debugging information"),

		GitProvider:        flag.String("vcs", "", "Specify version control system to scan (Eg. github, gitlab, bitbucket)"),
		GitProviderBaseURL: flag.String("baseurl", "", "Specify VCS base URL"),
		GitProviderToken:   flag.String("token", "", "Specify VCS token"),
		EnvFilePath:        flag.String("env", "", ".env file path containing VCS base URLs and tokens"),
		RepoID:             flag.String("repo-id", "", "Scan the repository with this ID"),
		ScanTargets:        flag.String("scan-targets", "", "Comma separated list of sub-directories within the repository to scan"),
		Repos:              flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
		GitScanPath:        flag.String("git-scan-path", "", "Specify the local path to scan"),
	}

	flag.Parse()

	return options, nil
}
