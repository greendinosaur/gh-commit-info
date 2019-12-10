package repos

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
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
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
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
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
		Err: nil,
	})

	GetRepoPRs(c)

	assert.EqualValues(t, http.StatusOK, response.Code)

	result := []github.MultiplePullRequestResponse{}
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
	//need to add in the parameters

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myowner/myrepo/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
		Err: nil,
	})

	GetRepoPRs(c)

	assert.EqualValues(t, http.StatusOK, response.Code)

	result := []github.MultiplePullRequestResponse{}
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
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
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
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`)),
		},
		Err: nil,
	})

	GetRepoSinglePR(c)

	assert.EqualValues(t, http.StatusOK, response.Code)

	var result github.GetSinglePullRequestResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", result.URL)
	assert.EqualValues(t, 123456, result.ID)
	assert.EqualValues(t, 9, result.Number)
	assert.EqualValues(t, "open", result.State)
	assert.EqualValues(t, "Title of the PR", result.Title)
}

//TODO: missing tests for the three new functions
