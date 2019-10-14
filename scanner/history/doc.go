package history

type ScanHistory struct {
	GitProvider string `json:"git_provider"`
	RepoID      string `json:"repo_id"`
	CommitHash  string `json:"commit_hash"`
	CreatedAt   string `json:"created_at"`
}
