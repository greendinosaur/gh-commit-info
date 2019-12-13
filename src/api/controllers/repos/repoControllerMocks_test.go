//mock the service layer so controller calls the mock and not the actual service layer
package repos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/domain/githubdomain"
	"github.com/greendinosaur/gh-commit-info/src/api/services"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	funcGetRepoPRs          func(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError)
	funcGetRepoSinglePR     func(owner string, repo string, pullRequst string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError)
	funcGetSingleCommitPR   func(owner string, repo string, SHA string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError)
	funcGetRepoCommits      func(owner string, repo string) ([]githubdomain.GetCommitInfo, errors.APIError)
	funcGetRepoSingleCommit func(owner string, repo string, SHA string) (*githubdomain.GetCommitInfo, errors.APIError)
	funcGetCodeReviewReport func(owner string, repo string, fromDate time.Time, endDate time.Time) (string, errors.APIError)
)

type repoServiceMock struct{}

func (s *repoServiceMock) GetRepoPRs(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {
	return funcGetRepoPRs(owner, repo, scope)
}

func (s *repoServiceMock) GetRepoSinglePR(owner string, repo string, pullRequest string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError) {
	return funcGetRepoSinglePR(owner, repo, pullRequest)
}

func (s *repoServiceMock) GetSingleCommitPR(owner string, repo string, SHA string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {
	return funcGetSingleCommitPR(owner, repo, SHA)
}

func (s *repoServiceMock) GetRepoCommits(owner string, repo string) ([]githubdomain.GetCommitInfo, errors.APIError) {
	return funcGetRepoCommits(owner, repo)
}

func (s *repoServiceMock) GetRepoSingleCommit(owner string, repo string, SHA string) (*githubdomain.GetCommitInfo, errors.APIError) {
	return funcGetRepoSingleCommit(owner, repo, SHA)
}

func (s *repoServiceMock) GetCodeReviewReport(owner string, repo string, fromDate time.Time, endDate time.Time) (string, errors.APIError) {
	return funcGetCodeReviewReport(owner, repo, fromDate, endDate)
}

func TestGetPRsNoErrorMockingEntireService(t *testing.T) {
	services.RepositoryService = &repoServiceMock{}

	funcGetRepoPRs = func(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {
		repoBase := githubdomain.RepoBase{
			Label: "A label",
			Ref:   "A Reference",
			SHA:   "ABCDEF123456768",
		}

		gitUser1 := githubdomain.GitUser{
			Login:     "My Login ID",
			ID:        123456,
			Type:      "A user",
			SiteAdmin: true,
		}

		gitUser2 := githubdomain.GitUser{
			Login:     "A Second Login ID",
			ID:        8767,
			Type:      "A user",
			SiteAdmin: false,
		}
		getMultiplePullRequestResponse := githubdomain.GetSinglePullRequestResponse{
			URL:            "some URL",
			ID:             123456,
			Number:         9,
			State:          "open",
			Title:          "Title of the PR",
			CreatedAt:      time.Now().AddDate(0, 0, -1).UTC(),
			UpdatedAt:      time.Now().AddDate(0, -1, 0).UTC(),
			ClosedAt:       time.Now().AddDate(0, -1, 0).UTC(),
			MergedAt:       time.Now().AddDate(0, -1, 0).UTC(),
			MergeCommitSHA: "ABCDEF1234567890",
			User:           gitUser1,
			Assignee:       gitUser2,
			Base:           repoBase,
		}
		result1 := []githubdomain.GetSinglePullRequestResponse{getMultiplePullRequestResponse}
		return result1, nil
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls?state=all", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	GetRepoPRs(c)

	assert.EqualValues(t, http.StatusOK, response.Code)

	result := []githubdomain.GetSinglePullRequestResponse{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", result[0].URL)
	assert.EqualValues(t, 123456, result[0].ID)
	assert.EqualValues(t, 9, result[0].Number)
	assert.EqualValues(t, "open", result[0].State)
	assert.EqualValues(t, "Title of the PR", result[0].Title)

}

func TestGetPRGithubErrorMockingEntireService(t *testing.T) {
	services.RepositoryService = &repoServiceMock{}

	funcGetRepoPRs = func(owner string, repo string, scope string) ([]githubdomain.GetSinglePullRequestResponse, errors.APIError) {

		return nil, errors.NewBadRequestError("invalid owner parameter")
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myrepo/pulls?state=all", strings.NewReader(`{}`))
	params := map[string]string{"repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	GetRepoPRs(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	APIErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, APIErr)
	assert.EqualValues(t, http.StatusBadRequest, APIErr.Status())
	assert.EqualValues(t, "invalid owner parameter", APIErr.Message())

}

func TestRepoGetSinglePRNoErrorMockingEntireService(t *testing.T) {
	services.RepositoryService = &repoServiceMock{}

	funcGetRepoSinglePR = func(owner string, repo string, pullRequest string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError) {
		repoBase := githubdomain.RepoBase{
			Label: "A label",
			Ref:   "A Reference",
			SHA:   "ABCDEF123456768",
		}

		gitUser1 := githubdomain.GitUser{
			Login:     "My Login ID",
			ID:        123456,
			Type:      "A user",
			SiteAdmin: true,
		}

		gitUser2 := githubdomain.GitUser{
			Login:     "A Second Login ID",
			ID:        8767,
			Type:      "A user",
			SiteAdmin: false,
		}
		getPRInfoResponse := githubdomain.GetSinglePullRequestResponse{
			URL:            "some URL",
			ID:             123456,
			Number:         9,
			State:          "open",
			Title:          "Title of the PR",
			CreatedAt:      time.Now().AddDate(0, 0, -1).UTC(),
			UpdatedAt:      time.Now().AddDate(0, -1, 0).UTC(),
			ClosedAt:       time.Now().AddDate(0, -1, 0).UTC(),
			MergedAt:       time.Now().AddDate(0, -1, 0).UTC(),
			MergeCommitSHA: "ABCDEF1234567890",
			User:           gitUser1,
			Assignee:       gitUser2,
			Base:           repoBase,
		}
		return &getPRInfoResponse, nil
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls/1", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "pull": "1"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	GetRepoSinglePR(c)

	assert.EqualValues(t, http.StatusOK, response.Code)

	var result githubdomain.GetSinglePullRequestResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", result.URL)
	assert.EqualValues(t, 123456, result.ID)
	assert.EqualValues(t, 9, result.Number)
	assert.EqualValues(t, "open", result.State)
	assert.EqualValues(t, "Title of the PR", result.Title)

}

func TestGetRepoSinglePRGithubErrorMockingEntireService(t *testing.T) {
	services.RepositoryService = &repoServiceMock{}

	funcGetRepoSinglePR = func(owner string, repo string, pullRequest string) (*githubdomain.GetSinglePullRequestResponse, errors.APIError) {

		return nil, errors.NewBadRequestError("invalid pull parameter")
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myrepo/pulls?state=all", strings.NewReader(`{}`))
	params := map[string]string{"owner": "owner", "repo": "myrepo", "pull": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	GetRepoSinglePR(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	APIErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, APIErr)
	assert.EqualValues(t, http.StatusBadRequest, APIErr.Status())
	assert.EqualValues(t, "invalid pull parameter", APIErr.Message())

}
