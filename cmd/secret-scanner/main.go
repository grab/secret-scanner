package main

import (
	"fmt"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/remotegit/github"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/remotegit/gitlab"
	"os"
	"time"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/logic/scan"
)

func main() {
	options, err := scan.ParseOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch *options.Source {
	case "github":
		sess, err := github.NewGithubSession(options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sess.Out.Important("Github Scanning Started at %s\n", sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures\n", len(scan.Signatures))
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
	case "gitlab":
		sess, err := gitlab.NewGitlabSession(options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sess.Out.Important("Gitlab Scanning Started at %s\n", sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures\n", len(scan.Signatures))
		if sess.Stats.Status == "finished" {
			sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
		} else {
			gitlab.GatherGitlabRepos(sess)
			gitlab.AnalyzeGitlabRepositories(sess)
			sess.Finish()
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
	default:
		fmt.Println("Specify version control system to scan (Eg. github, gitlab, bitbucket)")
		os.Exit(1)
	}
}
