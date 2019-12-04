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
	c.JSON(http.StatusCreated, result)
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
	c.JSON(http.StatusCreated, result)
}
