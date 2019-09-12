package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/common/signatures"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/remotegit/github"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/options"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
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

	if *opt.EnvFilePath != "" {
		err = godotenv.Load(*opt.EnvFilePath)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error loading .env file path %s: %v", *opt.EnvFilePath, err))
			os.Exit(1)
		}
	}

	var gitProvider gitprovider.GitProvider

	// Initialize new scan session
	sess := &session.Session{}
	sess.Initialize(opt)
	sess.Out.Important("%s Scanning Started at %s\n", strings.Title(*opt.GitProvider), sess.Stats.StartedAt.Format(time.RFC3339))
	sess.Out.Important("Loaded %d signatures\n", len(signatures.Signatures))

	// Set Git provider
	switch *opt.GitProvider {
	case "github":
		sess, err := github.NewGithubSession(opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if sess.Stats.Status == "finished" {
			sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
		} else {
			if len(sess.Options.Logins) == 0 {
				sess.Out.Fatal("Please provide at least one GitHub organization or user\n")
			}

			github.GatherTargets(sess)
			github.GatherRepositories(sess)
			github.AnalyzeRepositories(sess)
			sess.Finish()
			sess.Out.Important("Github Scanning Finished at %s\n", sess.Stats.FinishedAt.Format(time.RFC3339))
			if *sess.Options.Save != "" {
				err := sess.SaveToFile(*sess.Options.Save)
				if err != nil {
					sess.Out.Error("Error saving session to %s: %s\n", *sess.Options.Save, err)
				}
				sess.Out.Important("Saved session to: %s\n\n", *sess.Options.Save)
			}
			sess.Stats.PrintStats(sess.Out)
		}
		return
	case "gitlab":
		gitProvider = &gitprovider.GitlabProvider{}
	default:
		fmt.Println("Specify Git provider type (Currently supports github, gitlab)")
		os.Exit(1)
	}

	// Initialize Git provider
	err = gitProvider.Initialize(os.Getenv("GITLAB_BASE_URL"), os.Getenv("GITLAB_TOKEN"), nil)
	if err != nil {
		sess.Out.Fatal("%v", err)
		os.Exit(1)
	}

	// Scan
	if sess.Stats.Status == "finished" {
		sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
		return
	}
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
