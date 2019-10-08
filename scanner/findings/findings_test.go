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
		CommitURL:       "",
		RepositoryURL:   "",
	}
}
