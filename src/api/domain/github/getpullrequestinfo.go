package github

import "time"

//GetSinglePullRequestResponse stores information about a single PR
type GetSinglePullRequestResponse struct {
	URL               string    `json:"url"`
	ID                int64     `json:"id"`
	Number            int64     `json:"number"`
	State             string    `json:"state"`
	Title             string    `json:"title"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ClosedAt          time.Time `json:"closed_at"`
	MergedAt          time.Time `json:"merged_at"`
	MergeCommitSHA    string    `json:"merge_commit_sha"`
	User              GitUser   `json:"user"`
	Assignee          GitUser   `json:"assignee"`
	Base              RepoBase  `json:"base"`
	AuthorAssociation string    `json:"author_association"`
	Draft             bool      `json:"draft"`
	Merged            bool      `json:"merged"`
	MergeableState    string    `json:"mergeable_state"`
	MergedBy          GitUser   `json:"merged_by"`
	Commits           int64     `json:"commits"`
}

//RepoBase stores info about the base of the repo
type RepoBase struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
}

//MultiplePullRequestResponse represents multiple pull requests
type MultiplePullRequestResponse struct {
	URL            string    `json:"url"`
	ID             int       `json:"id"`
	Number         int       `json:"number"`
	State          string    `json:"state"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ClosedAt       time.Time `json:"closed_at"`
	MergedAt       time.Time `json:"merged_at"`
	MergeCommitSHA string    `json:"merge_commit_sha"`
	User           GitUser   `json:"user"`
	Assignee       GitUser   `json:"assignee"`
	Base           RepoBase  `json:"base"`
}
