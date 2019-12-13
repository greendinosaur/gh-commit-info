package githubdomain

import "time"

//GetCommitInfo returns information about a commit
type GetCommitInfo struct {
	URL           string                        `json:"url"`
	SHA           string                        `json:"sha"`
	Commit        DetailedCommitInfo            `json:"commit"`
	Author        GitUser                       `json:"author"`
	Committer     GitUser                       `json:"committer"`
	Parents       []Parent                      `json:"parents"`
	IsMergeCommit bool                          `json:"ismergecommit"` //not set by github, calculated later in code
	PRForMerge    *GetSinglePullRequestResponse `json:"pull"`          //not set by github
}

//DetailedCommitInfo has more detailed info about the commit
type DetailedCommitInfo struct {
	URL       string     `json:"url"`
	Author    CommitUser `json:"author"`
	Committer CommitUser `json:"committer"`
	Message   string     `json:"message"`
}

//CommitUser has info about the user doing the commit
type CommitUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

//Parent has information about the parent of the commit
type Parent struct {
	URL string `json:"url"`
	SHA string `json:"sha"`
}
