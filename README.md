# Secret Scanner

Secret scanner is a command-line tool to scan Git repositories for any sensitive information such as private keys, API secrets and tokens, etc.

It does so by looking at file names, extensions, and content, attempting to match them against a list of signatures.

The too is loosely based on <a href="https://github.com/michenriksen/gitrob">Gitrob</a>, with added support for Gitlab on top of Github.

For more information: https://wiki.grab.com/display/IS/Code+secrets+scanner

## Setup

The use of this tool requires you to set VCS (Github / Gitlab / Bitbucket) API base URL and tokens in your environment.

You can do so by:
```
export GITHUB_BASE_URL=https://api.github.com; GITHUB_TOKEN=my-token
```

Alternatively, you can add them into your `.bash_profile` or create a `.env` file. See `.env.example`.

You can also provide them as tool flag options. See **Usage** section

The precedence of usage is as follows from highest to lowest:
1. Flag options
2. .env file
3. CLI exported
4. Values from .bash_profile

## Usage

Basic:
```
secret-scanner -vcs github -env .env -repo-list repo.csv
```

CLI Args:
```
-baseurl string
    Specify VCS base URL

-commit-depth int
    Number of repository commits to process (default 500)

-debug
    Print debugging information

-env string
    .env file path containing VCS base URLs and tokens

-git-scan-path string
    Specify the local path to scan

-load string
    Load session file

-repo-id string
    Scan the repository with this ID

-repo-list string
    CSV file containing the list of whitelisted repositories to scan

-save string
    Save session to file

-scan-targets string
    Comma separated list of sub-directories within the repository to scan

-silent
    Suppress all output except for errors

-threads int
    Number of concurrent threads (default number of logical CPUs)

-token string
    Specify VCS token

-vcs string
    Specify version control system to scan (Eg. github, gitlab, bitbucket)
```
