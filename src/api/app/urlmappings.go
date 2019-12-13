package app

import (
	"github.com/greendinosaur/gh-commit-info/src/api/controllers/bobby"
	"github.com/greendinosaur/gh-commit-info/src/api/controllers/repos"
)

func mapURLs() {
	router.GET("/bobby", bobby.Chariot)
	router.GET("/repos/:owner/:repo/pulls", repos.GetRepoPRs)
	router.GET("/repos/:owner/:repo/pulls/:pull", repos.GetRepoSinglePR)
	router.GET("/repos/:owner/:repo/commits", repos.GetRepoCommits)
	router.GET("/repos/:owner/:repo/commits/:sha", repos.GetRepoSingleCommit)
	router.GET("/repos/:owner/:repo/commits/:sha/pulls", repos.GetPRsForSingleCommit)
	router.GET("/codereview/:owner/:repo", repos.GetCodeReviewReport)

}
