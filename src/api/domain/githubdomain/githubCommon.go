package githubdomain

//GitUser stores info about a user
//Is referenced in both Pull Requests and Commits
type GitUser struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}
