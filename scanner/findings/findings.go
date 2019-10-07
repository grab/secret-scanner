package findings

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Finding struct {
	Id              string
	FilePath        string
	Action          string
	Description     string
	Comment         string
	RepositoryOwner string
	RepositoryName  string
	CommitHash      string
	CommitMessage   string
	CommitAuthor    string
	FileUrl         string
	CommitUrl       string
	RepositoryUrl   string
}

func (f *Finding) generateID() {
	// Used for dedupe in defect dojo
	h := sha1.New()
	io.WriteString(h, f.FileUrl)
	io.WriteString(h, f.Action)

	// io.WriteString(h, f.CommitHash)
	// io.WriteString(h, f.CommitMessage)
	// io.WriteString(h, f.CommitAuthor)
	io.WriteString(h, f.Description)
	f.Id = fmt.Sprintf("%x", h.Sum(nil))
}

func (f *Finding) Initialize() {
	f.generateID()
}
