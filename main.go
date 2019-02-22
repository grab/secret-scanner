package main

import (
  "fmt"
  "os"
	"time"
	
  "./core"
)


func main() {
	options, err := core.ParseOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
  } 

	switch *options.Source {
		case "github":
			sess, err := core.NewGithubSession(options)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			sess.Out.Important("Github Scanning Started at %s\n", sess.Stats.StartedAt.Format(time.RFC3339))
			sess.Out.Important("Loaded %d signatures\n", len(core.Signatures))
			if sess.Stats.Status == "finished" {
				sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
			} else {
				if len(sess.Options.Logins) == 0 {
					sess.Out.Fatal("Please provide at least one GitHub organization or user\n")
				}
		
				core.GatherTargets(sess)
				core.GatherRepositories(sess)
				core.AnalyzeRepositories(sess)
				sess.Finish()
				sess.Out.Important("Github Scanning Finished at %s\n", sess.Stats.FinishedAt.Format(time.RFC3339))
				if *sess.Options.Save != "" {
					err := sess.SaveToFile(*sess.Options.Save)
					if err != nil {
						sess.Out.Error("Error saving session to %s: %s\n", *sess.Options.Save, err)
					}
					sess.Out.Important("Saved session to: %s\n\n", *sess.Options.Save)
				}
				core.PrintSessionStats(sess.Session)
			}
		case "gitlab":
			sess, err := core.NewGitlabSession(options)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			sess.Out.Important("Gitlab Scanning Started at %s\n", sess.Stats.StartedAt.Format(time.RFC3339))
			sess.Out.Important("Loaded %d signatures\n", len(core.Signatures))
			if sess.Stats.Status == "finished" {
				sess.Out.Important("Loaded session file: %s\n", *sess.Options.Load)
			}	else {
				core.GatherGitlabRepos(sess)
				core.AnalyzeGitlabRepositories(sess)
				sess.Finish()
				sess.Out.Important("Gitlab Scanning Finished at %s\n", sess.Stats.FinishedAt.Format(time.RFC3339))
				if *sess.Options.Save != "" {
					err := sess.SaveToFile(*sess.Options.Save)
					if err != nil {
						sess.Out.Error("Error saving session to %s: %s\n", *sess.Options.Save, err)
					}
					sess.Out.Important("Saved session to: %s\n\n", *sess.Options.Save)
				}
				core.PrintSessionStats(sess.Session)
			}
		default:
			fmt.Println("Specify version control system to scan (Eg. github, gitlab, bitbucket)")
			os.Exit(1)		
	}

	// sess.Out.Important("Press Ctrl+C to stop web server and exit.\n\n")
	// select {}
}

