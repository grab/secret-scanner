package options

import (
	"flag"
	"os"
	"strings"
)

type Options struct {
	CommitDepth       *int
	GithubAccessToken *string `json:"-"`
	NoExpandOrgs      *bool
	Threads           *int
	Save              *string `json:"-"`
	Load              *string `json:"-"`
	Silent            *bool
	Debug             *bool
	Logins            []string

	GitProvider       *string
	EnvFilePath       *string
	RepoID            *string
	ScanTargets       *string
	Repos             *string
}

func (o Options) ValidateOptions() error {
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
		CommitDepth:       flag.Int("commit-depth", 500, "Number of repository commits to process"),
		GithubAccessToken: flag.String("github-access-token", "", "GitHub access token to use for API requests"),
		NoExpandOrgs:      flag.Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations"),
		Threads:           flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Save:              flag.String("save", "", "Save session to file"),
		Load:              flag.String("load", "", "Load session file"),
		Silent:            flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:             flag.Bool("debug", false, "Print debugging information"),

		GitProvider:       flag.String("git-provider", "", "Specify version control system to scan (Eg. github, gitlab, bitbucket)"),
		EnvFilePath:       flag.String("env", "", ".env file path containing Git provider base URL and tokens"),
		RepoID:            flag.String("repo-id", "", "Scan the repository with this ID"),
		ScanTargets:       flag.String("scan-targets", "", "Comma separated list of sub-directories within the repository to scan"),
		Repos:             flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
	}

	flag.Parse()
	options.Logins = flag.Args()

	return options, nil
}
