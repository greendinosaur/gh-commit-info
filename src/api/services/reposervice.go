package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/config"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/githubdomain"
	"github.com/greendinosaur/gh-commit-info/src/api/providers/githubprovider"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	GetRepoPRs(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError)
	GetRepoSinglePR(owner string, repo string, pullNumber string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError)
	GetSingleCommitPR(owner string, repo string, SHA string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError)
	GetRepoCommits(owner string, repo string) ([]githubdomain.GetCommitInfo, errors.APIError)
	GetRepoSingleCommit(owner string, repo string, SHA string) (*githubdomain.GetCommitInfo, errors.APIError)
	GetCodeReviewReport(owner string, repo string, fromDate time.Time, endDate time.Time) (string, errors.APIError)
}

const (
	errorInvalidOwnerParam = "invalid owner parameter"
	errorInvalidRepoParam  = "invalid repo parameter"
	errorInvalidScopeParam = "invalid scope parameter"
	errorInvalidSHAParam   = "invalid SHA parameter"
	errorInvalidPullParam  = "invalid pull parameter"
)

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
		return owner, repo, scope, errors.NewBadRequestError(errorInvalidOwnerParam)
	}

	if len(repo) == 0 {
		return owner, repo, scope, errors.NewBadRequestError(errorInvalidRepoParam)
	}
	if len(scope) == 0 {
		return owner, repo, scope, errors.NewBadRequestError(errorInvalidScopeParam)
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
	return owner, repo, scope, errors.NewBadRequestError(errorInvalidScopeParam)

}

//check that the inputs to get a single PR are good
func validateSinglePRInputs(owner string, repo string, pullNumber string) (string, string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	pullNumber = strings.TrimSpace(pullNumber)

	if len(owner) == 0 {
		//an error
		return owner, repo, pullNumber, errors.NewBadRequestError(errorInvalidOwnerParam)
	}

	if len(repo) == 0 {
		return owner, repo, pullNumber, errors.NewBadRequestError(errorInvalidRepoParam)
	}
	if len(pullNumber) == 0 {
		return owner, repo, pullNumber, errors.NewBadRequestError(errorInvalidPullParam)
	}

	if _, err := strconv.ParseInt(pullNumber, 10, 64); err != nil {
		return owner, repo, pullNumber, errors.NewBadRequestError(errorInvalidPullParam)

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
		return owner, repo, SHA, errors.NewBadRequestError(errorInvalidOwnerParam)
	}

	if len(repo) == 0 {
		return owner, repo, SHA, errors.NewBadRequestError(errorInvalidRepoParam)
	}
	if len(SHA) == 0 {
		return owner, repo, SHA, errors.NewBadRequestError(errorInvalidSHAParam)
	}

	return owner, repo, SHA, nil

}

//check that the inputs for all commits in a repo are good
func validateAllCommitsInputs(owner string, repo string) (string, string, errors.APIError) {

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)

	if len(owner) == 0 {
		//an error
		return owner, repo, errors.NewBadRequestError(errorInvalidOwnerParam)
	}

	if len(repo) == 0 {
		return owner, repo, errors.NewBadRequestError(errorInvalidRepoParam)
	}

	return owner, repo, nil

}

//GetRepoPRs returns pull request information for the given repo
func (s *reposService) GetRepoPRs(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {
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
func (s *reposService) GetRepoSinglePR(owner string, repo string, pullNumber string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError) {
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
func (s *reposService) GetSingleCommitPR(owner string, repo string, SHA string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {

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
func (s *reposService) GetRepoCommits(owner string, repo string) ([]githubdomain.GetCommitInfo, errors.APIError) {
	var err errors.APIError
	owner, repo, err = validateAllCommitsInputs(owner, repo)
	if err != nil {
		return nil, err
	}

	response, errProvider := githubprovider.GetRepoCommits(config.GetGithubAccessToken(), owner, repo)
	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//getRepoCommitsInDateRange returns all commits from the given repo in the indicated date range
func getRepoCommitsInDateRange(owner string, repo string, fromDate time.Time, toDate time.Time) ([]githubdomain.GetCommitInfo, errors.APIError) {
	var err errors.APIError
	owner, repo, err = validateAllCommitsInputs(owner, repo)
	if err != nil {
		return nil, err
	}

	response, errProvider := githubprovider.GetRepoCommitsInDateRange(config.GetGithubAccessToken(), owner, repo, fromDate, toDate)

	if errProvider != nil {
		return nil, errors.NewAPIError(errProvider.StatusCode, errProvider.Message)
	}

	return response, nil
}

//GetRepoSingleCommit returns details about a specific commit inside the indicated repo
func (s *reposService) GetRepoSingleCommit(owner string, repo string, SHA string) (*githubdomain.GetCommitInfo, errors.APIError) {
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

//isMergeCommit determines if a commit is a merge commit
func isMergeCommit(commitInfo *githubdomain.GetCommitInfo) bool {
	//business logic from github that a merge commit has two parents, other commits don't
	//assumption is that commits have been retrieved from the master branch only

	if len(commitInfo.Parents) > 1 {
		return true
	}
	return false

}

func isPRResultingInMerge(pullRequest *githubdomain.GetSinglePullRequestResponse) bool {
	if pullRequest.State == "closed" && len(pullRequest.MergeCommitSHA) > 0 {
		return true
	}

	return false
}

//GetCodeReviewReport returns a text file that summarises the commit and PR data and also
//provides a list of the relevant commits and PRs
//this function will need to:
//1. get all the commits in a given timeframe
//2. for each commit, determine if a merge commit or proper commit
//3. for each proper commit, determine if it has an approved PR that resulted in the merge commit
//4. summarise the results (#total commits, #merge commits, #commits with PR, #commits with no PR)
//5. summarise the commits (sha, committer, date, commit message)
//6. summarise the PRs (PR title, approver, raiser, date)
//7. output this all as a text file that can be streamed back via the a client via an API
//Note PR reviews are stored in a different object and require a different GitHub API call
//therefore, if want to get the list of approvers who approved the PR, need another API call
//should look to do this
func (s *reposService) GetCodeReviewReport(owner string, repo string, fromDate time.Time, endDate time.Time) (string, errors.APIError) {

	//set-up the counters for the statistics to report on later
	totalMergeCommits := 0
	totalCommitsWithPR := 0
	totalCommitsWithNoPR := 0

	//create an array where we can store indices for those commits that aren't merge commits and have no PRs
	var indexCommitsWithNoPR []int

	repoCommits, err := getRepoCommitsInDateRange(owner, repo, fromDate, endDate)

	if err != nil {
		return "", err
	}

	//now we can loop over each commit and get hold of the associated PRs
	for commitCounter, repoCommitInfo := range repoCommits {
		//check for the merge and store this
		isMergeCommit := isMergeCommit(&repoCommitInfo)
		repoCommitInfo.IsMergeCommit = isMergeCommit

		if isMergeCommit {
			totalMergeCommits++
		}

		//now get the associated PRs and find one that has been closed and has a merge commit
		//may be multiple PRs associated with this commit
		//still making assumption feature branches are created, then a PR and merge into master
		//so simple feature branching with development on master
		//cater for other scenarios later
		pullsForCommit, err := RepositoryService.GetSingleCommitPR(owner, repo, repoCommitInfo.SHA)

		if err != nil {
			return "", err
		}

		if len(pullsForCommit) == 0 {
			//no PR associated with this commit so need to store for reporting purposes
			indexCommitsWithNoPR = append(indexCommitsWithNoPR, commitCounter)
			totalCommitsWithNoPR++
		} else {
			//now we have the array of PRs iterate through and see if there is a merged and closed PR
			for _, pull := range pullsForCommit {
				if isPRResultingInMerge(&pull) {
					//put details of this PR into the commit for later reporting
					repoCommitInfo.PRForMerge = &pull
					totalCommitsWithPR++
					//TODO: what if multiple PRs? does this happen?
				}
			}
		}

	}

	//can now summarise the information and return this
	summaryInfo := fmt.Sprintf("#Total Commits: %d, #Merged Commits: %d,  #Commits with PRs: %d, #Non-merge commits with no PR: %d",
		len(repoCommits), totalMergeCommits, totalCommitsWithPR, totalCommitsWithNoPR)

	return summaryInfo, nil
}
