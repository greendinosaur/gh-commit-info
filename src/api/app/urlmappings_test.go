package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	mapURLs()
	os.Exit(m.Run())
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {

	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w

}

func TestMapURLsServiceRunning(t *testing.T) {

	w := performRequest(router, "GET", "/bobby")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "bobby and his chariots", w.Body.String())
}

func TestMapURLsNoMap(t *testing.T) {
	w := performRequest(router, "GET", "/bobbyfgfg")
	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestGetRepoPRsErrorFromGithub(t *testing.T) {

	gin.SetMode(gin.TestMode)

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

	w := performRequest(router, "GET", "/repos/myowner/myrepo/pulls?state=all")

	assert.EqualValues(t, http.StatusUnauthorized, w.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(w.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestGetRepoSinglePRErrorFromGithub(t *testing.T) {

	gin.SetMode(gin.TestMode)

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

	w := performRequest(router, "GET", "/repos/myowner/myrepo/pulls/1")

	assert.EqualValues(t, http.StatusUnauthorized, w.Code)
	apiErr, err := errors.NewAPIErrorFromBytes(w.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

//TODO: missing tests for the three commit based API calls, not added as likely these won't be exposed
//in this way in the future as the API exposed to clients will change
