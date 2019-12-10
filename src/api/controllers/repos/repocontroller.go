package repos

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greendinosaur/gh-commit-info/src/api/services"
)

//GetRepoPRs returns the pull requests for the given repo
func GetRepoPRs(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	state := c.Query("state")
	fmt.Println(state)

	if state == "" {
		state = "all"
	}
	fmt.Println(state)
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
