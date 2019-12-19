/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package findings

import (
	"crypto/sha256"
	"fmt"
	"io"
)

const (
	// MaxLineChar defines the maximum number of characters in line content
	MaxLineChar = 100
)

// Finding holds the info for scan finding
type Finding struct {
	ID              string
	FilePath        string
	Action          string
	Description     string
	Comment         string
	RepositoryOwner string
	RepositoryName  string
	CommitHash      string
	CommitMessage   string
	CommitAuthor    string
	FileURL         string
	Line            uint64
	LineContent     string
	CommitURL       string
	RepositoryURL   string
	IsTestContext   bool
}

// GenerateHashID generates an unique hash
func (f *Finding) GenerateHashID() (hash string, err error) {
	// Used for dedupe in defect dojo
	h := sha256.New()
	str := fmt.Sprintf("%s%s%v%s", f.FileURL, f.Action, f.Line, f.LineContent)

	_, err = io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

	// io.WriteString(h, f.CommitHash)
	// io.WriteString(h, f.CommitMessage)
	// io.WriteString(h, f.CommitAuthor)
}

// TruncateLineContent truncates line content
func (f *Finding) TruncateLineContent(maxLen int) {
	lineLen := len(f.LineContent)
	if maxLen > 0 && lineLen > 0 && lineLen > maxLen {
		f.LineContent = f.LineContent[0:maxLen]
	}
}
