//Package githubprovider provides commit and PR information from Github
package githubprovider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
)

//information needed to get commit data from Github
const (
	urlGetRepoCommits      = "https://api.github.com/repos/%s/%s/commits"
	urlGetRepoSingleCommit = "https://api.github.com/repos/%s/%s/commits/%s"
)

//GetRepoCommits returns commits for the given repo
func GetRepoCommits(accessToken string, owner string, repo string) ([]github.GetCommitInfo, *github.GithubErrorResponse) {
	URL := fmt.Sprintf(urlGetRepoCommits, owner, repo)
	headers := getCommonHeader(accessToken)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	var result []github.GetCommitInfo

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal successful response: %s", err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return result, nil
}

//GetRepoSingleCommit returns details about a single commit
func GetRepoSingleCommit(accessToken string, owner string, repo string, sha string) (*github.GetCommitInfo, *github.GithubErrorResponse) {

	URL := fmt.Sprintf(urlGetRepoSingleCommit, owner, repo, sha)
	headers := getCommonHeader(accessToken)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	var result github.GetCommitInfo

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal successful response: %s", err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return &result, nil
}
