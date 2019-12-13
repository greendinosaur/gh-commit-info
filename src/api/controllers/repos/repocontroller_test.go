package repos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/githubdomain"
	"github.com/greendinosaur/gh-commit-info/src/api/providers/githubprovider"
	"github.com/greendinosaur/gh-commit-info/src/api/services"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestGetRepoPRsErrorFromGithub(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls?state=all", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetRepoPRs(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetRepoPRsNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	services.ResetService()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls?state=all", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataPRsResponseStatusCode(),
			Body:       testutils.GetMockDataPRsResponseMessage(),
		},
		Err: nil,
	})

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

func TestGetRepoPRsMissingStateParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	services.ResetService()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}

	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataPRsResponseStatusCode(),
			Body:       testutils.GetMockDataPRsResponseMessage(),
		},
		Err: nil,
	})

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

func TestGetRepoSinglePRErrorFromGithub(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls/1", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "pull": "1"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetRepoSinglePR(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetRepoSinglePRNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	services.ResetService()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/pulls/1", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "pull": "1"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSinglePRResponseStatusCode(),
			Body:       testutils.GetMockDataSinglePRResponseMessage(),
		},
		Err: nil,
	})

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

//TODO: missing tests for the three new functions
func TestGetRepoCommitsWithError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetRepoCommits(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetRepoCommitsNoError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataCommitsResponseStatusCode(),
			Body:       testutils.GetMockDataCommitsResponseMessage(),
		},
		Err: nil,
	})

	GetRepoCommits(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	var result []githubdomain.GetCommitInfo
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, "http://www.github.com", result[0].URL)
	assert.EqualValues(t, "AABCDEF123456", result[0].SHA)
}

func TestGetRepoSingleCommitGithubError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/SHA123", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "sha": "SHA123"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits/SHA123",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetRepoSingleCommit(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetRepoSingleCommitNoError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/AABCDEF123456", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "sha": "AABCDEF123456"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits/AABCDEF123456",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataSingleCommitResponseMessage(),
		},
		Err: nil,
	})

	GetRepoSingleCommit(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	var result githubdomain.GetCommitInfo
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, "http://www.github.com", result.URL)
	assert.EqualValues(t, "AABCDEF123456", result.SHA)
}

func TestGetPRsForSingleCommitInvalidResponse(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/SHA123/pulls", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "sha": "SHA123"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits/SHA123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetPRsForSingleCommit(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetPRsForSingleCommitNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	services.ResetService()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/repos/myowner/myrepo/commits/SHA123/pulls", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo", "sha": "SHA123"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/commits/SHA123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataPRsResponseStatusCode(),
			Body:       testutils.GetMockDataPRsResponseMessage(),
		},
		Err: nil,
	})

	GetPRsForSingleCommit(c)

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

func TestGetCodeReviewReportError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/codereview/myowner/myrepo", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myowner", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myowner/myrepo/commits?since=" + fromDate.UTC().Format(githubprovider.FmtGithubDate) + "&until=" + toDate.UTC().Format(githubprovider.FmtGithubDate)

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
		Err: nil,
	})

	GetCodeReviewReport(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestCodeReviewReportNoError(t *testing.T) {
	services.ResetService()
	gin.SetMode(gin.TestMode)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/codereview/myuser/myrepo", strings.NewReader(`{}`))
	params := map[string]string{"owner": "myuser", "repo": "myrepo"}
	c, _ := testutils.GetMockedContextWithParams(request, response, params)

	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(githubprovider.FmtGithubDate) + "&until=" + toDate.UTC().Format(githubprovider.FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataSingleSliceNonMergeCommitResponsesMessage(),
		},
	})

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits/AABCDEF123456/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataApprovedPRForCommitResponsesMessage(),
		},
	})

	GetCodeReviewReport(c)

	result := string(response.Body.Bytes())
	assert.EqualValues(t, "#Total Commits: 1, #Merged Commits: 0,  #Commits with PRs: 1, #Commits with No PRs: 0", result)

}
