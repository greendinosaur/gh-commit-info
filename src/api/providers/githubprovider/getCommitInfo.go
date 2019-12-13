//Package githubprovider provides commit and PR information from Github
package githubprovider

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/domain/githubdomain"
)

//information needed to get commit data from Github
const (
	urlGetRepoCommits            = "https://api.github.com/repos/%s/%s/commits"
	urlGetRepoSingleCommit       = "https://api.github.com/repos/%s/%s/commits/%s"
	urlGetRepoCommitsInDateRange = "https://api.github.com/repos/%s/%s/commits?since=%s&until=%s"
)

//GetRepoCommits returns commits for the given repo
func GetRepoCommits(accessToken string, owner string, repo string) ([]githubdomain.GetCommitInfo, *githubdomain.GithubErrorResponse) {
	URL := fmt.Sprintf(urlGetRepoCommits, owner, repo)
	headers := getCommonHeader(accessToken)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	var result []githubdomain.GetCommitInfo

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf(errorUnmarshallingResponse, err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return result, nil
}

//GetRepoCommitsInDateRange returns commits for the given repo
func GetRepoCommitsInDateRange(accessToken string, owner string, repo string, fromDate time.Time, toDate time.Time) ([]githubdomain.GetCommitInfo, *githubdomain.GithubErrorResponse) {
	URL := fmt.Sprintf(urlGetRepoCommitsInDateRange, owner, repo, fromDate.UTC().Format(FmtGithubDate), toDate.UTC().Format(FmtGithubDate))
	headers := getCommonHeader(accessToken)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	var result []githubdomain.GetCommitInfo

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf(errorUnmarshallingResponse, err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return result, nil
}

//GetRepoSingleCommit returns details about a single commit
func GetRepoSingleCommit(accessToken string, owner string, repo string, sha string) (*githubdomain.GetCommitInfo, *githubdomain.GithubErrorResponse) {

	URL := fmt.Sprintf(urlGetRepoSingleCommit, owner, repo, sha)
	headers := getCommonHeader(accessToken)

	bytes, err := getDataFromGithubAPI(URL, headers)

	if err != nil {
		return nil, err
	}

	var result githubdomain.GetCommitInfo

	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf(errorUnmarshallingResponse, err.Error()))
		return nil, getUnmarshalBodyError()
	}
	return &result, nil
}
