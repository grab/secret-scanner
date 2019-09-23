package options

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"net/url"
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
	BaseURL            *string
	Token              *string
	//ClientID           *string
	//ClientSecret       *string
	//UserID             *string
	//UserPW             *string
	EnvFilePath        *string
	RepoID             *string
	ScanTarget         *string
	Repos              *string
	GitScanPath        *string
}

func (o Options) ValidateOptions() (bool, error) {
	if *o.RepoID != "" && *o.Repos != "" {
		return false, ErrRepoOptionConflict
	}
	if *o.EnvFilePath != "" {
		if _, err := os.Stat(*o.EnvFilePath); os.IsNotExist(err) {
			return false, err
		}
	}

	// Load env file if present
	if *o.EnvFilePath != "" {
		err := godotenv.Load(*o.EnvFilePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("error: unable to load .env file path %s: %v", *o.EnvFilePath, err))
			os.Exit(1)
		}
	}

	if *o.GitProvider != gitprovider.GithubName && *o.GitProvider != gitprovider.GitlabName && *o.GitProvider != gitprovider.BitbucketName {
		return false, ErrInvalidGitProvider
	}

	switch *o.GitProvider {
	case gitprovider.GithubName:
		return o.ValidateGithubOptions(), nil
	case gitprovider.GitlabName:
		return o.ValidateGitlabOptions(), nil
	case gitprovider.BitbucketName:
		return o.ValidateBitbucketOptions(), nil
	default:
		return false, ErrInvalidGitProvider
	}
}

func (o Options) ValidateGithubOptions() bool {
	return true
}

func (o Options) ValidateGitlabOptions() bool {
	baseURL := *o.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("GITLAB_BASE_URL")
		*o.BaseURL = baseURL
	}
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return false
	}
	return o.ValidateHasToken("GITLAB_TOKEN")
}

func (o Options) ValidateBitbucketOptions() bool {
	return true
}

func (o Options) ValidateHasToken(key string) bool {
	if *o.Token == "" {
		if os.Getenv(key) == "" {
			return false
		}
		*o.Token = os.Getenv(key)
	}
	return true
}

func (o *Options) ParseScanTargets() []string {
	return strings.Split(*o.ScanTarget, ",")
}

func Parse() (Options, error) {
	options := Options{
		CommitDepth:        flag.Int("commit-depth", 500, "Number of repository commits to process"),
		Threads:            flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Save:               flag.String("save", "", "Save session to file"),
		Load:               flag.String("load", "", "Load session file"),
		Silent:             flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:              flag.Bool("debug", false, "Print debugging information"),

		GitProvider:        flag.String("git", "", "Specify type of git provider (Eg. github, gitlab, bitbucket)"),
		BaseURL:            flag.String("baseurl", "", "Specify VCS base URL"),
		Token:              flag.String("token", "", "Specify VCS token"),
		//ClientID:           flag.String("oauth-id", "", "Specify Bitbucket Oauth2 client ID"),
		//ClientSecret:       flag.String("oauth-secret", "", "Specify Bitbucket Oauth2 client secret"),
		//UserID:             flag.String("user-id", "", "Specify Bitbucket username"),
		//UserPW:             flag.String("user-pw", "", "Specify Bitbucket password"),
		EnvFilePath:        flag.String("env", "", ".env file path containing VCS base URLs and tokens"),
		RepoID:             flag.String("repo-id", "", "Scan the repository with this ID"),
		ScanTarget:         flag.String("scan-target", "", "Sub-directory within the repository to scan"),
		Repos:              flag.String("repo-list", "", "CSV file containing the list of whitelisted repositories to scan"),
		GitScanPath:        flag.String("git-scan-path", "", "Specify the local path to scan"),
	}

	flag.Parse()

	return options, nil
}
