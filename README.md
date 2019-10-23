# Secret Scanner

Secret scanner is a command-line tool to scan Git repositories for any sensitive information such as private keys, API secrets and tokens, etc.

It does so by looking at file names, extensions, and content, attempting to match them against a list of signatures.

The tool is based on <a href="https://github.com/michenriksen/gitrob">Gitrob</a>, with added support for Gitlab and Bitbucket on top of Github.

## Setup

## Auth Tokens

The use of this tool requires you to set various Git provider (Github / Gitlab / Bitbucket) authentication token in your environment.

You can do so by:
```
export GITHUB_TOKEN=my-token; secret-scanner -repos jquery/jquery
```

To persist the various Git provider tokens, you can add them into your `.bash_profile` or create a `.env` file. See `.env.example`.

### Skip Files

You can define paths to be excluded from scanning by defining them in a comma separated format in `.env` file.

`SKIP_EXT` defines the file extensions to be excluded
`SKIP_PATHS` defines the paths/files to be excluded if the path matches one of the patterns defined in the list
`SKIP_TEST_PATHS` defines any test directories/files that you would like to skip. It is being kept separately from `SKIP_PATHS` because sometimes it may be useful to scan the test files as well. You can toggle to scan test files by giving `-skip-tests=false` in the CLI.

## Usage

For `bool` CLI flags, use `=` between key-val pair. Eg `-ui=false`

### Basic

The most basic usage requires a list of Github repository identifiers in the form of `org/repo`.

```
./secret-scanner -repos jquery/jquery
```

To scan repositories in other Git providers, simply specify the Git provider name.

For Gitlab, provide the project ID instead of `org/repo`.

```
./secret-scanner -git bitbucket -repos litmis/mama
./secret-scanner -git gitlab -repos 3836952
```

You can scan multiple repositories from the same Git provider by providing multiple identifiers separated by commas.

```
./secret-scanner -repos jquery/jquery,lodash/lodash
```

### Local Scan

By default, the tool will attempt to make a clone before scanning the files.

If you already have a copy of the repository on local disk, you can do a local scan by specifying the `dir` parameter.

```
./secret-scanner -dir /dir/path/to/local/repository
```

### Sub-directory Scan

In instances where a repository contains multiple projects (i.e monorepo), or you simply want to scan specific sub-directory, you can do so by providing `sub-dir`.

Example:
https://github.com/user/awesome-projects contains
- build/
- dist/
- src/
- test/
- ...

Caveat: Only works with single `repos`

To scan `src` only:
```
./secret-scanner -repos jquery/jquery -sub-dir src
```

## Scan Results as Output

By default, findings found during the scan will be printed as console output. You can save it as JSON to path by specifying the `output` param

```
./secret-scanner -repos jquery/jquery -output ~/report.json
```

The output file will contain the lines containing the potential secrets. In circumstances where you do not want to expose them, you can specify `-log-secret=false`

## Scan State

By default, no scan state is being kept, meaning every scan on the same repository will start afresh.

If scan state is enabled, the scanner will save the latest scan session and commit hash in JSON format. From the next scan onwards for the same repository,the scanner will only scan changes since the last saved commit hash.

The default location of scan state JSON file is in `~/.secret-scanner/`.

```
./secret-scanner -repos jquery/jquery -use-state=true
```

## CLI Args

```
  -baseurl string
        Specify Git provider base URL

  -commit-depth int
        Number of repository commits to process (default 500)

  -debug
        Print debugging information

  -env string
        .env file path containing Git provider base URLs and tokens

  -git string
        Name of git provider (Eg. github, gitlab, bitbucket) (default "github")

  -load string
        Load session file

  -dir string
        Specify the local git repo path to scan

  -log-secret
        If true, the matched secret will be included in output file (default true)

  -output string
        Save session to file

  -repos string
        Comma-separated list of repos to scan

  -sub-dir string
        Sub-directory within the repository to scan

  -quiet
        Suppress all output except for errors

  -skip-tests
        Skips possible test contexts (default true)

  -use-state
        If use-state is off, every scan will be treated as a brand new scan.

  -threads int
        Number of concurrent threads (default number of logical CPUs)

  -token string
        Specify Git provider token
```

## Credits

Project is built upon the ground work laid in <a href="https://github.com/michenriksen/gitrob" target="_blank">Gitrob</a> by <a href="https://michenriksen.com/" target="_blank">Michael Henriksen</a>.

And many secret signatures was taken from <a href="https://github.com/eth0izzle/shhgit/" target="_blank">shhgit</a>.
