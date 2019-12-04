package services

import (
	"log"
	"strconv"
	"strings"

	"github.com/greendinosaur/gh-commit-info/src/api/config"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
	"github.com/greendinosaur/gh-commit-info/src/api/providers/githubprovider"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	GetRepoPRs(owner string, repo string, scope string) ([]github.MultiplePullRequestResponse, errors.APIError)
	GetRepoSinglePR(owner string, repo string, pullNumber string) (*github.GetSinglePullRequestResponse, errors.APIError)
}

//RepositoryService defines the service to use
var RepositoryService reposServiceInterface

func init() {
	RepositoryService = &reposService{}
}

//ResetService calls the init function again
func ResetService() {
	RepositoryService = &reposService{}
}

func validateInputs(owner string, repo string, scope string) (string, string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	scope = strings.TrimSpace(scope)

	if len(owner) == 0 {
		//an error
		return owner, repo, scope, errors.NewBadRequestError("invalid owner parameter")
	}

	if len(repo) == 0 {
		return owner, repo, scope, errors.NewBadRequestError("invalid repo parameter")
	}
	if len(scope) == 0 {
		return owner, repo, scope, errors.NewBadRequestError("invalid scope parameter")
	}

	switch scope {
	case
		"open",
		"closed",
		"all":
		//if got this far then owner, repo and scope are all okay
		return owner, repo, scope, nil
	}
	//if this far then scope is bad
	return owner, repo, scope, errors.NewBadRequestError("invalid scope parameter")

}

func validateSinglePRInputs(owner string, repo string, pullNumber string) (string, string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	pullNumber = strings.TrimSpace(pullNumber)

	if len(owner) == 0 {
		//an error
		return owner, repo, pullNumber, errors.NewBadRequestError("invalid owner parameter")
	}

	if len(repo) == 0 {
		return owner, repo, pullNumber, errors.NewBadRequestError("invalid repo parameter")
	}
	if len(pullNumber) == 0 {
		return owner, repo, pullNumber, errors.NewBadRequestError("invalid pull parameter")
	}

	if _, err := strconv.ParseInt(pullNumber, 10, 64); err != nil {
		return owner, repo, pullNumber, errors.NewBadRequestError("invalid pull parameter")

	}

	return owner, repo, pullNumber, nil

}

//GetRepoPRs returns pull request information for the given repo
func (s *reposService) GetRepoPRs(owner string, repo string, scope string) ([]github.MultiplePullRequestResponse, errors.APIError) {
	//firstly check the input params are valid and create an error otherwise
	var err errors.APIError
	owner, repo, scope, err = validateInputs(owner, repo, scope)
	if err != nil {
		return nil, err
	}
	//then call the provider with valid parameters

	response, errProvider := githubprovider.GetRepoPRs(config.GetGithubAccessToken(), owner, repo, scope)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//GetRepoSinglePR returns details about a single pull request
func (s *reposService) GetRepoSinglePR(owner string, repo string, pullNumber string) (*github.GetSinglePullRequestResponse, errors.APIError) {
	//firstly check the input params are valid and create an error otherwise
	var err errors.APIError
	owner, repo, pullNumber, err = validateSinglePRInputs(owner, repo, pullNumber)
	if err != nil {
		return nil, err
	}
	//then call the provider with valid parameters
	log.Println("calling provider", owner, repo, pullNumber)
	response, errProvider := githubprovider.GetRepoSinglePR(config.GetGithubAccessToken(), owner, repo, pullNumber)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}
