# Secret Scanner

Secret scanner is a command-line tool to scan Git repositories for any sensitive information such as private keys, API secrets and tokens, etc.

It does so by looking at file names, extensions, and content, attempting to match them against a list of signatures.

The tool is based on <a href="https://github.com/michenriksen/gitrob">Gitrob</a>, with added support for Gitlab and Bitbucket on top of Github.

For more information: https://wiki.grab.com/display/IS/Code+secrets+scanner

## Setup

The use of this tool requires you to set various Git provider (Github / Gitlab / Bitbucket) options such as API base URL and tokens in your environment.

You can do so by:
```
export GITHUB_BASE_URL=https://api.github.com; GITHUB_TOKEN=my-token
```

Alternatively, you can add them into your `.bash_profile` or create a `.env` file. See `.env.example`.

You can also provide them as tool flag options. See **Usage** section

The precedence of usage is as follows from highest to lowest:
1. Command flag options
2. .env file
3. CLI exported
4. Values from .bash_profile

## Usage

For `bool` cmd-line flags, use `=` between key-val pair. Eg `-ui=false`

### Basic

```
./secret-scanner -git github -env .env -repo-list repo.csv
./secret-scanner -git bitbucket -env .env -repo-list repo.csv
./secret-scanner -git gitlab -baseurl https://mygitlab.com -token mysecret-token -repo-list repo.csv
```

### Local Scan

By default, the tool will attempt to make a local clone before scanning the files.

If you already have a copy of the files on local disk, you can do a local scan by specifying the `git-scan-path` parameter.

```
./secret-scanner -git github -git-scan-path /path/to/local/repository
```

### Sub-directory Scan

In instances where a repository contains multiple projects (i.e monorepo), you can specify which project to scan by providing `scan-target`, the project directory path name relative to repository root.

Example:
https://github.com/user/awesome-projects contains
- dir1/
- dir2/
- dir3/

To scan `dir1`:
```
./secret-scanner -git github -env .env -repo-list repo.csv -git-scan-path /path/to/awesome-projects -scan-target dir1
```

## Report

By default, findings found during the scan will be printed as console output. You can save it as JSON by specifying the `save` param

```
./secret-scanner -git github -env .env -repo-list repo.csv -save ./report.json
```

### Web UI

By default, the after the scan is completed, a local web server will be spun up containing the findings in a nice UI.

You can turn it off by specifying `ui` to false.

```
./secret-scanner -git github -env .env -repo-list repo.csv -save ./report.json -ui false
```

## CLI Args

```
-baseurl string
     Specify Git provider base URL

-commit-depth int
    Number of repository commits to process (default 500)

-debug bool
    Print debugging information

-env string
    .env file path containing Git provider base URLs and tokens

-git string
    Specify type of git provider (Eg. github, gitlab, bitbucket)

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

-scan-target string
    Sub-directory within the repository to scan

-silent bool
    Suppress all output except for errors

-skip-tests bool
     Skips possible test contexts

-threads int
    Number of concurrent threads (default number of logical CPUs)

-token string
    Specify VCS token

-ui bool
    Serves up local UI for scan results if true, (default true)
```

## Credits

Project is built upon the ground work laid in <a href="https://github.com/michenriksen/gitrob" target="_blank">Gitrob</a> by <a href="https://michenriksen.com/" target="_blank">Michael Henriksen</a>.

And many secret signatures was taken from <a href="https://github.com/eth0izzle/shhgit/" target="_blank">shhgit</a>.
