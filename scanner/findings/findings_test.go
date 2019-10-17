/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package findings

import "testing"

func TestFinding_GenerateHashID(t *testing.T) {
	f := createNewFinding()
	hashID, err := f.GenerateHashID()
	if err != nil {
		t.Errorf("Want no err, got err")
		return
	}
	if len(hashID) == 0 {
		t.Errorf("Want %v, got 0", len(hashID))
	}
}

func TestFinding_TruncateLineContent(t *testing.T) {
	finding := createNewFinding()
	finding.LineContent = "this is a line content with 47 characters in it"

	finding.TruncateLineContent(10)
	if finding.LineContent != "this is a " {
		t.Errorf("Want \"this is a \", got %v", finding.LineContent)
	}

	finding.TruncateLineContent(0)
	if finding.LineContent != "this is a " {
		t.Errorf("Want \"this is a \", got %v", finding.LineContent)
	}
}

func createNewFinding() *Finding {
	return &Finding{
		ID:              "",
		FilePath:        "",
		Action:          "",
		Description:     "",
		Comment:         "",
		RepositoryOwner: "",
		RepositoryName:  "",
		CommitHash:      "",
		CommitMessage:   "",
		CommitAuthor:    "",
		FileURL:         "",
		Line:            0,
		LineContent:     "",
		CommitURL:       "",
		RepositoryURL:   "",
		IsTestContext:   false,
	}
}
