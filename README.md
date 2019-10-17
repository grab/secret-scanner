# Secret Scanner

Secret scanner is a command-line tool to scan Git repositories for any sensitive information such as private keys, API secrets and tokens, etc.

It does so by looking at file names, extensions, and content, attempting to match them against a list of signatures.

The tool is based on <a href="https://github.com/michenriksen/gitrob">Gitrob</a>, with added support for Gitlab and Bitbucket on top of Github.

## Setup

The use of this tool requires you to set various Git provider (Github / Gitlab / Bitbucket) options such as API base URL, tokens, etc in your environment if you would like to scan your own private repositories.

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

### Skip Files

You can define paths to be excluded from scanning by defining them in a comma separated format in `.env` file.

`SKIP_EXT` defines the file extensions to be excluded
`SKIP_PATHS` defines the paths/files to be excluded if the path matches one of the patterns defined in the list
`SKIP_TEST_PATHS` defines any test directories/files that you would like to skip. It is being kept separately from `SKIP_PATHS` because sometimes it may be useful to scan the test files as well. You can toggle to scan test files by giving `-skip-test=false` in the CLI.

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

## Scan History

By default, no scan history is being kept, meaning every scan on the same repository will start afresh.

If scan history is enabled, the scanner will save the latest scan session and commit hash in JSON format. From the next scan onwards for the same repository,the scanner will only scan changes since the last saved commit hash.

The default location of scan history JSON file is in `~/.secret-scanner/`. You can define the place you want to store the history by giving `-history my/custom/path/to/history` in the CLI

```
./secret-scanner -git github -env .env -repo-list repo.csv -save ./report.json -no-history=false -history my/custom/path/to/history
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

-history string
    File path to store scan histories

-load string
    Load session file

-log-secret bool
    If true, the matched secret will be included in results save file (default true)

-no-history bool
    If no-history is on, every scan will be treated as a brand new scan. (default true)

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
