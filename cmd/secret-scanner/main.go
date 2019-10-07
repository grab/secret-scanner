package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/web"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/gitprovider"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/options"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/signatures"
)

func main() {
	opt, err := options.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate Options
	optValid, err := opt.ValidateOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !optValid {
		fmt.Println(errors.New("invalid option(s)"))
		os.Exit(1)
	}

	var gitProvider gitprovider.GitProvider
	additionalParams := map[string]string{}

	// Set Git provider
	switch *opt.GitProvider {
	case gitprovider.GithubName:
		gitProvider = &gitprovider.GithubProvider{}

	case gitprovider.GitlabName:
		gitProvider = &gitprovider.GitlabProvider{}

	case gitprovider.BitbucketName:
		gitProvider = &gitprovider.BitbucketProvider{}

	default:
		fmt.Println("error: invalid Git provider type (Currently supports github, gitlab)")
		os.Exit(1)
	}

	// Initialize new scan session
	sess := &session.Session{}
	sess.Initialize(opt)
	sess.Out.Important("%s Scanning Started at %s\n", strings.Title(*opt.GitProvider), sess.Stats.StartedAt.Format(time.RFC3339))
	sess.Out.Important("Loaded %d signatures\n", len(signatures.Signatures))

	// Initialize Git provider
	err = gitProvider.Initialize(*sess.Options.BaseURL, *sess.Options.Token, additionalParams)
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

	// Serve UI
	if *sess.Options.UI == "true" {
		web.InitRouter("127.0.0.1", "8888", sess)
	}
}
