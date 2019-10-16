/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package findings

import (
	"crypto/sha1"
	"fmt"
	"io"
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
	CommitURL       string
	RepositoryURL   string
	IsTestContext   bool
}

// GenerateHashID generates an unique hash
func (f *Finding) GenerateHashID() (hash string, err error) {
	// Used for dedupe in defect dojo
	h := sha1.New()
	str := fmt.Sprintf("%s%s%s", f.FileURL, f.Action, f.Description)

	_, err = io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

	// io.WriteString(h, f.CommitHash)
	// io.WriteString(h, f.CommitMessage)
	// io.WriteString(h, f.CommitAuthor)
}
