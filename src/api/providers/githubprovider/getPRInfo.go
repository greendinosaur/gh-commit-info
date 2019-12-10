package githubprovider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
)

//information needed to get PR data from Github
const (
	urlGetRepoPRs          = "https://api.github.com/repos/%s/%s/pulls?state=%s"
	urlGetRepoSinglePR     = "https://api.github.com/repos/%s/%s/pulls/%s"
	urlGetRepoPRForCommits = "https://api.github.com/repos/%s/%s/commits/%s/pulls"
)

//GetRepoSinglePR returns the given PR for a repo
func GetRepoSinglePR(accessToken string, owner string, repo string, pullNumber string) (*github.GetSinglePullRequestResponse, *github.GithubErrorResponse) {

	//setup the end point to call including the headers
	URL := fmt.Sprintf(urlGetRepoSinglePR, owner, repo, pullNumber)

	headers := getCommonHeader(accessToken)
	headers.Set(headerAccept, headerPRDraftAPI)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	//now we have a response, need to unmarshal it and return the results
	var result github.GetSinglePullRequestResponse

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal successful response: %s", err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return &result, nil
}

//getRepoPRsFromURL is used to return more than one pull request
func getRepoPRsFromURL(URL string, headers http.Header) ([]github.MultiplePullRequestResponse, *github.GithubErrorResponse) {

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	//now we have the response, unmarshal it back into the correct object to return
	var result []github.MultiplePullRequestResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal successful PR response: %s", err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return result, nil

}

//GetRepoPRs returns all of the PRs in the given repo
func GetRepoPRs(accessToken string, owner string, repo string, state string) ([]github.MultiplePullRequestResponse, *github.GithubErrorResponse) {

	//need to construct the URL to call and also the headers to send
	//these vary depending on the API call being made as described in the github API documentation
	URL := fmt.Sprintf(urlGetRepoPRs, owner, repo, state)

	headers := getCommonHeader(accessToken)
	headers.Set(headerAccept, headerPRDraftAPI)

	return getRepoPRsFromURL(URL, headers)
}

//GetSingleCommitPR returns all the PRs associated with a single commit SHA
func GetSingleCommitPR(accessToken string, owner string, repo string, SHA string) ([]github.MultiplePullRequestResponse, *github.GithubErrorResponse) {

	//construct the URL and headers
	//these can vary depending on the end point being called
	URL := fmt.Sprintf(urlGetRepoPRForCommits, owner, repo, SHA)

	headers := getCommonHeader(accessToken)
	headers.Set(headerAccept, headerPRForCommitDraftAPI)

	//can now call the API to get hold of the PRs
	return getRepoPRsFromURL(URL, headers)

}
