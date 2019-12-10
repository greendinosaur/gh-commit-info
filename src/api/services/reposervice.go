package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/config"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
	"github.com/greendinosaur/gh-commit-info/src/api/providers/githubprovider"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	GetRepoPRs(owner string, repo string, scope string) ([]github.MultiplePullRequestResponse, errors.APIError)
	GetRepoSinglePR(owner string, repo string, pullNumber string) (*github.GetSinglePullRequestResponse, errors.APIError)
	GetSingleCommitPR(owner string, repo string, SHA string) ([]github.MultiplePullRequestResponse, errors.APIError)
	GetRepoCommits(owner string, repo string) ([]github.GetCommitInfo, errors.APIError)
	GetRepoSingleCommit(owner string, repo string, SHA string) (*github.GetCommitInfo, errors.APIError)
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

//check the inputs for an individual PR are correct
func validatePRInputs(owner string, repo string, scope string) (string, string, string, errors.APIError) {

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

//check that the inputs to get a single PR are good
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

//check that the inputs for the PRs associated to a single commit are good
func validateSingleCommitPRInputs(owner string, repo string, SHA string) (string, string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	SHA = strings.TrimSpace(SHA)

	if len(owner) == 0 {
		//an error
		return owner, repo, SHA, errors.NewBadRequestError("invalid owner parameter")
	}

	if len(repo) == 0 {
		return owner, repo, SHA, errors.NewBadRequestError("invalid repo parameter")
	}
	if len(SHA) == 0 {
		return owner, repo, SHA, errors.NewBadRequestError("invalid SHA parameter")
	}

	return owner, repo, SHA, nil

}

//check that the inputs for all commits in a repo are good
func validateAllommitsInputs(owner string, repo string) (string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)

	if len(owner) == 0 {
		//an error
		return owner, repo, errors.NewBadRequestError("invalid owner parameter")
	}

	if len(repo) == 0 {
		return owner, repo, errors.NewBadRequestError("invalid repo parameter")
	}

	return owner, repo, nil

}

//GetRepoPRs returns pull request information for the given repo
func (s *reposService) GetRepoPRs(owner string, repo string, scope string) ([]github.MultiplePullRequestResponse, errors.APIError) {
	//firstly check the input params are valid and create an error otherwise
	var err errors.APIError
	owner, repo, scope, err = validatePRInputs(owner, repo, scope)
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
	response, errProvider := githubprovider.GetRepoSinglePR(config.GetGithubAccessToken(), owner, repo, pullNumber)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//GetSingleCommitPR returns the PRs associated with the specific commit SHA
func (s *reposService) GetSingleCommitPR(owner string, repo string, SHA string) ([]github.MultiplePullRequestResponse, errors.APIError) {

	var err errors.APIError
	owner, repo, SHA, err = validateSingleCommitPRInputs(owner, repo, SHA)
	if err != nil {
		return nil, err
	}

	response, errProvider := githubprovider.GetSingleCommitPR(config.GetGithubAccessToken(), owner, repo, SHA)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil

}

//GetRepoCommits returns all commits from the given repo
func (s *reposService) GetRepoCommits(owner string, repo string) ([]github.GetCommitInfo, errors.APIError) {
	var err errors.APIError
	owner, repo, err = validateAllommitsInputs(owner, repo)
	if err != nil {
		return nil, err
	}

	response, errProvider := githubprovider.GetRepoCommits(config.GetGithubAccessToken(), owner, repo)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//GetRepoSingleCommit returns details about a specific commit inside the indicated repo
func (s *reposService) GetRepoSingleCommit(owner string, repo string, SHA string) (*github.GetCommitInfo, errors.APIError) {
	var err errors.APIError
	owner, repo, SHA, err = validateSingleCommitPRInputs(owner, repo, SHA)
	if err != nil {
		return nil, err
	}

	response, errProvider := githubprovider.GetRepoSingleCommit(config.GetGithubAccessToken(), owner, repo, SHA)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//determines if a commit is a merge commit
func isMergeCommit() bool {
	return false
}

//determines if all commits in a given timeframe are associated with an approved PR
func isCommitWithApprovedPR(startTime time.Time, endtime time.Time) bool {
	return false
}

//checks that the person approving the PR is not responsible for any of the commits detailed in the PR
func isPRApproverDifferentToCommitters() bool {
	return false
}
