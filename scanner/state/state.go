/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package state

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// History contains scan history fields
type History struct {
	ID          string `json:"id"`
	GitProvider string `json:"git_provider"`
	RepoID      string `json:"repo_id"`
	CommitHash  string `json:"commit_hash"`
	CreatedAt   string `json:"created_at"`
}

// GetMapKey returns the scan history map key
func (h *History) GetMapKey() string {
	return fmt.Sprintf("%s:%s", h.GitProvider, h.RepoID)
}

// AssignID generates a hash and assign it to history ID field
func (h *History) AssignID() {
	hasher := sha256.New()

	_, err := io.WriteString(hasher, fmt.Sprintf("%s:%s", h.GitProvider, h.RepoID))
	if err == nil {
		h.ID = fmt.Sprintf("%x", hasher.Sum(nil))
	}
}

// Create creates a scan history struct
func Create(gitProvider, repoID, commitHash, createdAt string) *History {
	h := &History{
		GitProvider: gitProvider,
		RepoID:      repoID,
		CommitHash:  commitHash,
		CreatedAt:   createdAt,
	}

	h.AssignID()

	return h
}
