package app

import (
	"github.com/greendinosaur/gh-commit-info/src/api/controllers/bobby"
	"github.com/greendinosaur/gh-commit-info/src/api/controllers/repos"
)

func mapURLs() {
	router.GET("/bobby", bobby.Chariot)
	router.GET("/repos/:owner/:repo/pulls", repos.GetRepoPRs)
	router.GET("/repos/:owner/:repo/pulls/:pull", repos.GetRepoSinglePR)
}
