package signatures

import "regexp"

type PatternSignature struct {
	part        string
	match       *regexp.Regexp
	description string
	comment     string
}

func (s PatternSignature) Match(file MatchFile) bool {
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
		return false
	}

	return s.match.MatchString(*haystack)
}

func (s PatternSignature) Description() string {
	return s.description
}

func (s PatternSignature) Comment() string {
	return s.comment
}

func (s PatternSignature) Part() string {
	return s.part
}

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
		match:       regexp.MustCompile(`(xox(p-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+|b-[a-z0-9]+-[a-zA-Z0-9]+|a-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+-[a-z0-9]+))`),
		description: "Slack Token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(eyj[a-z0-9\\-_%]+.eyj[a-z0-9\\-_%]+.[a-z0-9\\-_%]+)|(refresh_token[\"']?\\s*[:=]\\s*[\"']?(?:[a-z0-9_]+-)+[a-z0-9_]+[\"']?)`),
		description: "Contains OAuth token",
		comment:     "",
	},
	PatternSignature{
		part:        PartContent,
		match:       regexp.MustCompile(`(aws|access|key|secret).*(([=:])|(:=))\\s*[\\\"']([A-Za-z0-9\/+=]{40})[\\\"']`),
		description: "AWS key",
		comment:     "",
	},
	//Creating lots of false positives
	// PatternSignature{
	// 	part:        PartPath,
	// 	match:       regexp.MustCompile(`credential`),
	// 	description: "Contains word: credential",
	// 	comment:     "",
	// },
	// PatternSignature{
	// 	part:        PartPath,
	// 	match:       regexp.MustCompile(`password`),
	// 	description: "Contains word: password",
	// 	comment:     "",
	// },
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`[g|G][i|I][t|T][l|L][a|A][b|B].*.[a-zA-Z0-9]{20}`),
	//	description: "Gitlab Token",
	//	comment:     "",
	//},
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`[g|G][i|I][t|T][h|H][u|U][b|B].*[0-9a-zA-Z]{35,40}`),
	//	description: "GitHub Token",
	//	comment:     "",
	//},
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*.[0-9a-zA-Z]{35,44}`),
	//	description: "Twitter Oauth 2",
	//	comment:     "",
	//},
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*.[0-9a-f]{32}`),
	//	description: "Facebook Oauth 2",
	//	comment:     "",
	//},
	//PatternSignature{
	//	part:        PartContent,
	//	match:       regexp.MustCompile(`[c|C][l|L][i|I][e|E][n|N][T|T][_][s|S][e|E][c|C][r|R][e|E][t|T].*[:].*[a-zA-Z0-9-_]{24}`),
	//	description: "Google Oauth 2",
	//	comment:     "",
	//},
}
