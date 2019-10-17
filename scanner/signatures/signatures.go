/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package signatures

import (
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
	for _, skippableExt := range skippableExtensions {
		if ext == skippableExt {
			return true
		}
	}
	for _, skippablePathIndicator := range skippablePathIndicators {
		if strings.Contains(path, skippablePathIndicator) {
			return true
		}
	}
	return false
}

// IsTestContext checks if file is in a test context
func (f *MatchFile) IsTestContext() bool {
	path := strings.ToLower(f.Path)

	for _, skippableTestContext := range skippableTestContexts {
		if strings.Contains(path, skippableTestContext) {
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
