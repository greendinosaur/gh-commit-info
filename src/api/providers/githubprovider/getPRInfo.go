package githubprovider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
)

const (
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"

	headerAccept     = "Accept"
	headerPRDraftAPI = "application/vnd.github.shadow-cat-preview+json"

	urlGetRepoCommits      = ""
	urlGetRepoPRs          = "https://api.github.com/repos/%s/%s/pulls?state=%s"
	urlGetRepoSinglePR     = "https://api.github.com/repos/%s/%s/pulls/%s"
	urlGetRepoSingleCommit = ""
	urlGetRepoCommitsForPR = ""
	urlGetRepoPRForCommit  = ""
)

func getAuthorizationHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func getPRsURL(owner string, repo string, state string) string {
	return fmt.Sprintf(urlGetRepoPRs, owner, repo, state)
}

func getSinglePRURL(owner string, repo string, pullNumber string) string {
	return fmt.Sprintf(urlGetRepoSinglePR, owner, repo, pullNumber)
}

//GetRepoSinglePR returns the given PR for a repo
func GetRepoSinglePR(accessToken string, owner string, repo string, pullNumber string) (*github.GetSinglePullRequestResponse, *github.GithubErrorResponse) {

	URL := getSinglePRURL(owner, repo, pullNumber)

	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(accessToken))
	headers.Set(headerAccept, headerPRDraftAPI)

	log.Println(URL, headers)

	response, err := restclient.Get(URL, headers)
	if err != nil {
		log.Println(fmt.Sprintf("error when getting pull requests: %s", err.Error()))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid response body"}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid json response body"}
		}

		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	var result github.GetSinglePullRequestResponse
	log.Println(string(bytes))
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal successful response: %s", err.Error()))
		log.Println("error in unmarshalling")
		log.Println("bytes:", string(bytes))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "error when trying to unmarshal github response"}
	}
	return &result, nil
}

//GetRepoPRs returns all of the PRs in the given repo
//TODO github will return an array so need to change the receiving structure to accept an array of Pull Requests and not a single one
func GetRepoPRs(accessToken string, owner string, repo string, state string) ([]github.MultiplePullRequestResponse, *github.GithubErrorResponse) {

	URL := getPRsURL(owner, repo, state)

	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(accessToken))
	headers.Set(headerAccept, headerPRDraftAPI)
	log.Println(URL)
	log.Println(headers)

	response, err := restclient.Get(URL, headers)
	if err != nil {
		log.Println(fmt.Sprintf("error when getting pull requests: %s", err.Error()))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid response body"}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid json response body"}
		}

		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	log.Println("about to unmarshal the multiple PRs response")
	log.Println(string(bytes))
	result := []github.MultiplePullRequestResponse{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal get pull requests successful response: %s", err.Error()))
		log.Println("error in unmarshalling")
		log.Println("bytes:", string(bytes))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "error when trying to unmarshal github response"}
	}
	log.Println(result)
	return result, nil
}

//func GetRepoPRsForCommit() {

//}

//put the commit related requests into a separate file, keep this one for pure PR data

//func GetRepoCommmits() {
//	return
//}

//func GetRepoSingleCommit() {

//}

//func GetRepoCommitsForPR() {

//}
