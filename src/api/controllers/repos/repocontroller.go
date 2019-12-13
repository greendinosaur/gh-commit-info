package repos

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greendinosaur/gh-commit-info/src/api/services"
)

//GetRepoPRs returns the pull requests for the given repo
func GetRepoPRs(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	state := c.Query("state")

	if state == "" {
		state = "all"
	}

	result, err := services.RepositoryService.GetRepoPRs(owner, repo, state)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetRepoSinglePR returns the indicated pull request in the repo
func GetRepoSinglePR(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	pullRequest := c.Param("pull")

	log.Println(owner, repo, pullRequest)

	result, err := services.RepositoryService.GetRepoSinglePR(owner, repo, pullRequest)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetRepoCommits returns all commits for a repo
func GetRepoCommits(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	result, err := services.RepositoryService.GetRepoCommits(owner, repo)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetRepoSingleCommit returns a single commit for a repo
func GetRepoSingleCommit(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	SHA := c.Param("sha")

	result, err := services.RepositoryService.GetRepoSingleCommit(owner, repo, SHA)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetPRsForSingleCommit returns the PRs associated to a specific commit
func GetPRsForSingleCommit(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	SHA := c.Param("sha")

	result, err := services.RepositoryService.GetSingleCommitPR(owner, repo, SHA)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetCodeReviewReport returns a plain text response with details of the commits and PRs
func GetCodeReviewReport(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	//TODO: for now hard code to the last month of commits, will need to pass this in as variables
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()

	result, err := services.RepositoryService.GetCodeReviewReport(owner, repo, fromDate, toDate)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte(result))
}
