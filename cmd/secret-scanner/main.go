package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/options"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/signatures"
	"os"
	"strings"
	"time"
)

func main() {
	opt, err := options.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate Options
	err = opt.ValidateOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Load env file if present
	if *opt.EnvFilePath != "" {
		err = godotenv.Load(*opt.EnvFilePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("error: unable to load .env file path %s: %v", *opt.EnvFilePath, err))
			os.Exit(1)
		}
	}

	var gitProvider gitprovider.GitProvider
	var gitProviderBaseURL string
	var gitProviderToken string

	// Set Git provider
	switch *opt.GitProvider {
	case "github":
		gitProvider = &gitprovider.GithubProvider{}
		gitProviderBaseURL = *opt.GitProviderBaseURL
		if gitProviderBaseURL == "" {
			gitProviderBaseURL = os.Getenv("GITHUB_BASE_URL")
		}
		gitProviderToken = *opt.GitProviderToken
		if gitProviderToken == "" {
			gitProviderToken = os.Getenv("GITHUB_TOKEN")
		}

	case "gitlab":
		gitProvider = &gitprovider.GitlabProvider{}
		gitProviderBaseURL = *opt.GitProviderBaseURL
		if gitProviderBaseURL == "" {
			gitProviderBaseURL = os.Getenv("GITLAB_BASE_URL")
		}
		gitProviderToken = *opt.GitProviderToken
		if gitProviderToken == "" {
			gitProviderToken = os.Getenv("GITLAB_TOKEN")
		}

	default:
		fmt.Println("error: invalid Git provider type (Currently supports github, gitlab)")
		os.Exit(1)
	}

	if gitProviderBaseURL == "" || gitProviderToken == "" {
		gitProviderUpper := strings.ToUpper(*opt.GitProvider)
		fmt.Println(fmt.Sprintf("error: VCS base URL and token not set. To set: \"export %s_BASE_URL=http://base-url.com; %s_TOKEN=my-token\" or set them in .env file", gitProviderUpper, gitProviderUpper))
		os.Exit(1)
	}

	// Initialize new scan session
	sess := &session.Session{}
	sess.Initialize(opt)
	sess.Out.Important("%s Scanning Started at %s\n", strings.Title(*opt.GitProvider), sess.Stats.StartedAt.Format(time.RFC3339))
	sess.Out.Important("Loaded %d signatures\n", len(signatures.Signatures))

	// Initialize Git provider
	err = gitProvider.Initialize(gitProviderBaseURL, gitProviderToken, nil)
	if err != nil {
		sess.Out.Fatal("%v", err)
		os.Exit(1)
	}

	if sess.Stats.Status == "finished" {
		sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
		return
	}

	// Scan
	scanner.Scan(sess, gitProvider)
	sess.Out.Important("Gitlab Scanning Finished at %s\n", sess.Stats.FinishedAt.Format(time.RFC3339))

	if *sess.Options.Save != "" {
		err := sess.SaveToFile(*sess.Options.Save)
		if err != nil {
			sess.Out.Error("Error saving session to %s: %s\n", *sess.Options.Save, err)
		}
		sess.Out.Important("Saved session to: %s\n\n", *sess.Options.Save)
	}

	sess.Stats.PrintStats(sess.Out)
}
