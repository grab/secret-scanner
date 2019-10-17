/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package signatures

// SimpleSignature ...
type SimpleSignature struct {
	part        string
	match       string
	description string
	comment     string
}

// Match checks if given file matches with signature
func (s SimpleSignature) Match(file MatchFile) []*MatchResult {
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

	if s.match == *haystack {
		matchResults = append(matchResults, &MatchResult{
			Filename:    file.Filename,
			Path:        file.Path,
			Extension:   file.Extension,
			Line:        0,
			LineContent: "",
		})
	}

	return matchResults
}

// Description returns signature description
func (s SimpleSignature) Description() string {
	return s.description
}

// Comment returns signature comment
func (s SimpleSignature) Comment() string {
	return s.comment
}

// Part returns signature part type
func (s SimpleSignature) Part() string {
	return s.part
}

// SimpleSignatures contains simple signatures
var SimpleSignatures = []Signature{
	// Extensions
	SimpleSignature{
		part:        PartExtension,
		match:       ".pem",
		description: "Potential cryptographic private key",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".log",
		description: "Log file",
		comment:     "Log files can contain secret HTTP endpoints, session IDs, API keys and other goodies",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".pkcs12",
		description: "Potential cryptographic key bundle",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".p12",
		description: "Potential cryptographic key bundle",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".pfx",
		description: "Potential cryptographic key bundle",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".asc",
		description: "Potential cryptographic key bundle",
		comment:     "",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "otr.private_key",
		description: "Pidgin OTR private key",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".ovpn",
		description: "OpenVPN client configuration file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".cscfg",
		description: "Azure service configuration schema file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".rdp",
		description: "Remote Desktop connection file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".mdf",
		description: "Microsoft SQL database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".sdf",
		description: "Microsoft SQL server compact database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".sqlite",
		description: "SQLite database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".sqlite3",
		description: "SQLite3 database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".bek",
		description: "Microsoft BitLocker recovery key file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".tpm",
		description: "Microsoft BitLocker Trusted Platform Module password file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".fve",
		description: "Windows BitLocker full volume encrypted data file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".jks",
		description: "Java keystore file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".psafe3",
		description: "Password Safe database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".agilekeychain",
		description: "1Password password manager database file",
		comment:     "Feed it to Hashcat and see if you're lucky",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".keychain",
		description: "Apple Keychain database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".pcap",
		description: "Network traffic capture file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".gnucash",
		description: "GnuCash database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".kwallet",
		description: "KDE Wallet Manager database file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".tblk",
		description: "Tunnelblick VPN configuration file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartExtension,
		match:       ".dayone",
		description: "Day One journal file",
		comment:     "Now it's getting creepy...",
	},

	// Filenames
	SimpleSignature{
		part:        PartFilename,
		match:       "secret_token.rb",
		description: "Ruby On Rails secret token configuration file",
		comment:     "If the Rails secret token is known, it can allow for remote code execution (http://www.exploit-db.com/exploits/27527/)",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "carrierwave.rb",
		description: "Carrierwave configuration file",
		comment:     "Can contain credentials for cloud storage systems such as Amazon S3 and Google Storage",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "database.yml",
		description: "Potential Ruby On Rails database configuration file",
		comment:     "Can contain database credentials",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "omniauth.rb",
		description: "OmniAuth configuration file",
		comment:     "The OmniAuth configuration file can contain client application secrets",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "settings.py",
		description: "Django configuration file",
		comment:     "Can contain database credentials, cloud storage system credentials, and other secrets",
	},

	SimpleSignature{
		part:        PartFilename,
		match:       "jenkins.plugins.publish_over_ssh.BapSshPublisherPlugin.xml",
		description: "Jenkins publish over SSH plugin file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "credentials.xml",
		description: "Potential Jenkins credentials file",
		comment:     "",
	},

	SimpleSignature{
		part:        PartFilename,
		match:       "LocalSettings.php",
		description: "Potential MediaWiki configuration file",
		comment:     "",
	},

	SimpleSignature{
		part:        PartFilename,
		match:       "Favorites.plist",
		description: "Sequel Pro MySQL database manager bookmark file",
		comment:     "",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "configuration.user.xpl",
		description: "Little Snitch firewall configuration file",
		comment:     "Contains traffic rules for applications",
	},

	SimpleSignature{
		part:        PartFilename,
		match:       "journal.txt",
		description: "Potential jrnl journal file",
		comment:     "Now it's getting creepy...",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "knife.rb",
		description: "Chef Knife configuration file",
		comment:     "Can contain references to Chef servers",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "proftpdpasswd",
		description: "cPanel backup ProFTPd credentials file",
		comment:     "Contains usernames and password hashes for FTP accounts",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "robomongo.json",
		description: "Robomongo MongoDB manager configuration file",
		comment:     "Can contain credentials for MongoDB databases",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "filezilla.xml",
		description: "FileZilla FTP configuration file",
		comment:     "Can contain credentials for FTP servers",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "recentservers.xml",
		description: "FileZilla FTP recent servers file",
		comment:     "Can contain credentials for FTP servers",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "ventrilo_srv.ini",
		description: "Ventrilo server configuration file",
		comment:     "Can contain passwords",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       "terraform.tfvars",
		description: "Terraform variable config file",
		comment:     "Can contain credentials for terraform providers",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       ".exports",
		description: "Shell configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       ".functions",
		description: "Shell configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
	SimpleSignature{
		part:        PartFilename,
		match:       ".extra",
		description: "Shell configuration file",
		comment:     "Shell configuration files can contain passwords, API keys, hostnames and other goodies",
	},
}
