/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package signatures

import (
	"regexp"
	"strings"
)

// PatternSignature ...
type PatternSignature struct {
	part        string
	match       *regexp.Regexp
	description string
	comment     string
}

// Match checks if given file matches with signature
func (s PatternSignature) Match(file MatchFile) []*MatchResult {
	var haystack *string
	switch s.part {
	case PartPath:
		haystack = &file.Path
	case PartFilename:
		haystack = &file.Filename
	case PartExtension:
		haystack = &file.Extension
	case PartContent:
		haystack = &file.Content
	default:
		return nil
	}

	var matchResults []*MatchResult
	contentBytes := []byte(file.ContentRaw)
	locations := s.match.FindAllIndex([]byte(*haystack), -1)

	for _, loc := range locations {
		contentBytesBefLine := contentBytes[0 : loc[1]-1]
		befLines := strings.Split(string(contentBytesBefLine), "\n")
		lineNo := len(befLines)

		matchResults = append(matchResults, &MatchResult{
			Filename:    file.Filename,
			Path:        file.Path,
			Extension:   file.Extension,
			Line:        uint64(lineNo),
			LineContent: string(contentBytes[loc[0]:loc[1]]),
		})
	}

	return matchResults
}

// Description returns signature description
func (s PatternSignature) Description() string {
	return s.description
}

// Comment returns signature comment
func (s PatternSignature) Comment() string {
	return s.comment
}

// Part returns signature part type
func (s PatternSignature) Part() string {
	return s.part
}

// PatternSignatures contains simple signatures
var PatternSignatures = []Signature{
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^.*_rsa$`),
		description: "Private SSH key",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^.*_dsa$`),
		description: "Private SSH key",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^.*_ed25519$`),
		description: "Private SSH key",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^.*_ecdsa$`),
		description: "Private SSH key",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?ssh/config$`),
		description: "SSH configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartExtension,
		match:       regexp.MustCompile(`^key(pair)?$`),
		description: "Potential cryptographic private key",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?(bash_|zsh_|sh_|z)?history$`),
		description: "Shell command history file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?mysql_history$`),
		description: "MySQL client command history file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?psql_history$`),
		description: "PostgreSQL client command history file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?pgpass$`),
		description: "PostgreSQL password file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?irb_history$`),
		description: "Ruby IRB console history file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?purple/accounts\.xml$`),
		description: "Pidgin chat client account configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?xchat2?/servlist_?\.conf$`),
		description: "Hexchat/XChat IRC client server list configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?irssi/config$`),
		description: "Irssi IRC client configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?recon-ng/keys\.db$`),
		description: "Recon-ng web reconnaissance framework API key database",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?dbeaver-data-sources.xml$`),
		description: "DBeaver SQL database manager configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?muttrc$`),
		description: "Mutt e-mail client configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?s3cfg$`),
		description: "S3cmd configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?aws/credentials$`),
		description: "AWS CLI credentials file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^sftp-config(\.json)?$`),
		description: "SFTP connection configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?trc$`),
		description: "T command-line Twitter client configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?gitrobrc$`),
		description: "Well, this is awkward... Gitrob configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?(bash|zsh|csh)rc$`),
		description: "Shell configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?(bash_|zsh_)?profile$`),
		description: "Shell profile configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?(bash_|zsh_)?aliases$`),
		description: "Shell command alias configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`config(\.inc)?\.php$`),
		description: "PHP configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartExtension,
		match:       regexp.MustCompile(`^key(store|ring)$`),
		description: "GNOME Keyring database file",
		comment:     "",
	},
	PatternSignature{
		part:        PartExtension,
		match:       regexp.MustCompile(`^kdbx?$`),
		description: "KeePass password manager database file",
		comment:     "Feed it to Hashcat and see if you're lucky",
	},
	PatternSignature{
		part:        PartExtension,
		match:       regexp.MustCompile(`^sql(dump)?$`),
		description: "SQL dump file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?htpasswd$`),
		description: "Apache htpasswd file",
		comment:     "",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^(\.|_)?netrc$`),
		description: "Configuration file for auto-login process",
		comment:     "Can contain username and password",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?gem/credentials$`),
		description: "Rubygems credentials file",
		comment:     "Can contain API key for a rubygems.org account",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?tugboat$`),
		description: "Tugboat DigitalOcean management tool configuration",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`doctl/config.yaml$`),
		description: "DigitalOcean doctl command-line client configuration file",
		comment:     "Contains DigitalOcean API key and other information",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?git-credentials$`),
		description: "git-credential-store helper credentials file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`config/hub$`),
		description: "GitHub Hub command-line client configuration file",
		comment:     "Can contain GitHub API access token",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?gitconfig$`),
		description: "Git configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`\.?chef/(.*)\.pem$`),
		description: "Chef private key",
		comment:     "Can be used to authenticate against Chef servers",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`etc/shadow$`),
		description: "Potential Linux shadow file",
		comment:     "Contains hashed passwords for system users",
	},
	PatternSignature{
		part:        PartPath,
		match:       regexp.MustCompile(`etc/passwd$`),
		description: "Potential Linux passwd file",
		comment:     "Contains system user information",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?dockercfg$`),
		description: "Docker configuration file",
		comment:     "Can contain credentials for public or private Docker registries",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?npmrc$`),
		description: "NPM configuration file",
		comment:     "Can contain credentials for NPM registries",
	},
	PatternSignature{
		part:        PartFilename,
		match:       regexp.MustCompile(`^\.?env$`),
		description: "Environment configuration file",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(?i)-{5}begin ([dr]sa|ec|openssh)? private key-{5}`),
		description: "Private Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|AGPA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`),
		description: "AWS Access Key ID Value",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)access)_?((?i)key)?_?((?i)id)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}(\\\"|'|`)?)"),
		description: "AWS Access Key ID",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)account)_?((?i)id)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[0-9]{4}-?[0-9]{4}-?[0-9]{4}(\\\"|'|`)?)"),
		description: "AWS Account ID",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)secret)_?((?i)access)?_?((?i)key)?_?((?i)id)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[A-Za-z0-9/+=]{40}(\\\"|'|`)?)"),
		description: "AWS Secret Access Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)session)?_?((?i)token)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[A-Za-z0-9/+=]{16,}(\\\"|'|`)?)"),
		description: "AWS Session Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("(?i)artifactory.{0,50}(\\\"|'|`)?[a-zA-Z0-9=]{112}(\\\"|'|`)?"),
		description: "Artifactory",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("(?i)codeclima.{0,50}(\\\"|'|`)?[0-9a-f]{64}(\\\"|'|`)?"),
		description: "CodeClimate",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`EAACEdEose0cBA[0-9A-Za-z]+`),
		description: "Facebook Access Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("((\\\"|'|`)?type(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?service_account(\\\"|'|`)?,?)"),
		description: "Google (GCM) Service account",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(?:r|s)k_[live|test]_[0-9a-zA-Z]{24}`),
		description: "Stripe API key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com`),
		description: "Google OAuth Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`AIza[0-9A-Za-z\\-_]{35}`),
		description: "Google Cloud API Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`ya29\\.[0-9A-Za-z\\-_]+`),
		description: "Google OAuth Access Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`sk_[live|test]_[0-9a-z]{32}`),
		description: "Picatic API key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`sq0atp-[0-9A-Za-z\-_]{22}`),
		description: "Square Access Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`sq0csp-[0-9A-Za-z\-_]{43}`),
		description: "Square OAuth Secret",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`access_token\$production\$[0-9a-z]{16}\$[0-9a-f]{32}`),
		description: "PayPal/Braintree Access Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`),
		description: "Amazon MWS Auth Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`SK[0-9a-fA-F]{32}`),
		description: "Twilo API Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`key-[0-9a-zA-Z]{32}`),
		description: "MailGun API Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`[0-9a-f]{32}-us[0-9]{12}`),
		description: "MailChimp API Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`sshpass -p.*['|\\\"]`),
		description: "SSH Password",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(https\\://outlook\\.office.com/webhook/[0-9a-f-]{36}\\@)`),
		description: "Outlook Team",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("(?i)sauce.{0,50}(\\\"|'|`)?[0-9a-f-]{36}(\\\"|'|`)?"),
		description: "Sauce Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(xox[pboa]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})`),
		description: "Slack Token",
		comment:     "",
	},
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`(xox(p-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+|b-[a-z0-9]+-[a-zA-Z0-9]+|a-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+))`),
	//	description: "Slack Token",
	//	comment:     "",
	//},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}`),
		description: "Slack Webhook",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("(?i)sonar.{0,50}(\\\"|'|`)?[0-9a-f]{40}(\\\"|'|`)?"),
		description: "SonarQube Docs API Key",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile("(?i)hockey.{0,50}(\\\"|'|`)?[0-9a-f]{32}(\\\"|'|`)?"),
		description: "HockeyApp",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`([\w+]{1,24})(://)([^$<]{1})([^\s";]{1,}):([^$<]{1})([^\s";]{1,})@[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,24}([^\s]+)`),
		description: "Username and password in URI",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(eyj[a-z0-9\\-_%]+.eyj[a-z0-9\\-_%]+.[a-z0-9\\-_%]+)|(refresh_token[\"']?\\s*[:=]\\s*[\"']?(?:[a-z0-9_]+-)+[a-z0-9_]+[\"']?)`),
		description: "Contains OAuth token",
		comment:     "",
	},
}
