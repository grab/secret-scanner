package signatures

import (
	"path/filepath"
	"strings"
)

type Signature interface {
	Match(file MatchFile) bool
	Description() string
	Comment() string
	Part() string
}

type MatchFile struct {
	Path      string
	Filename  string
	Extension string
	Content   string
}

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

func NewMatchFile(path string, content string) MatchFile {
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	content = strings.ToLower(content)
	return MatchFile{
		Path:      path,
		Filename:  filename,
		Extension: extension,
		Content:   content,
	}
}

func LoadSignatures() []Signature {
	sig := SimpleSignatures
	return append(sig, PatternSignatures...)
}
