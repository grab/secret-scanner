/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package signatures

import (
	"os"
	"path/filepath"
	"strings"
)

// Signature defines fields for a secret signature
type Signature interface {
	Match(file MatchFile) []*MatchResult
	Description() string
	Comment() string
	Part() string
}

// MatchFile contains details of a matching file
type MatchFile struct {
	Path       string
	Filename   string
	Extension  string
	Content    string
	ContentRaw string
}

// MatchResult contains match info
type MatchResult struct {
	Filename    string
	Path        string
	Extension   string
	Line        uint64
	LineContent string
}

// IsSkippable determines if a given matched file can be ignored
func (f *MatchFile) IsSkippable() bool {
	ext := strings.ToLower(f.Extension)
	path := strings.ToLower(f.Path)
	skipExts := skippableExtensions
	skipPaths := skippablePathIndicators

	if envSkipExt := os.Getenv("SKIP_EXT"); envSkipExt != "" {
		skipExts = strings.Split(envSkipExt, ",")
		for i, s := range skipExts {
			skipExts[i] = strings.TrimSpace(s)
		}
	}

	for _, skipExt := range skipExts {
		if ext == skipExt {
			return true
		}
	}

	if envSkipPaths := os.Getenv("SKIP_PATHS"); envSkipPaths != "" {
		skipPaths = strings.Split(envSkipPaths, ",")
		for i, s := range skipPaths {
			skipPaths[i] = strings.TrimSpace(s)
		}
	}

	for _, skipPath := range skipPaths {
		if strings.Contains(path, skipPath) {
			return true
		}
	}

	return false
}

// IsTestContext checks if file is in a test context
func (f *MatchFile) IsTestContext() bool {
	path := strings.ToLower(f.Path)
	skipTestPaths := skippableTestPaths

	if envSkipPaths := os.Getenv("SKIP_TEST_PATHS"); envSkipPaths != "" {
		skipTestPaths = strings.Split(envSkipPaths, ",")
		for i, s := range skipTestPaths {
			skipTestPaths[i] = strings.TrimSpace(s)
		}
	}

	for _, skipTestPath := range skipTestPaths {
		if strings.Contains(path, skipTestPath) {
			return true
		}
	}

	return false
}

// NewMatchFile creates new MatchFile
func NewMatchFile(path string, content string) MatchFile {
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	return MatchFile{
		Path:       path,
		Filename:   filename,
		Extension:  extension,
		Content:    strings.ToLower(content),
		ContentRaw: content,
	}
}

// LoadSignatures loads all signatures
func LoadSignatures() []Signature {
	sig := SimpleSignatures
	return append(sig, PatternSignatures...)
}
