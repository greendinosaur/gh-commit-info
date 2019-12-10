package githubprovider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuthorization)
	assert.EqualValues(t, "token %s", headerAuthorizationFormat)
	assert.EqualValues(t, "Accept", headerAccept)
	assert.EqualValues(t, "application/vnd.github.shadow-cat-preview+json", headerPRDraftAPI)
	assert.EqualValues(t, "application/vnd.github.groot-preview+json", headerPRForCommitDraftAPI)
	assert.EqualValues(t, "https://api.github.com/repos/%s/%s/pulls?state=%s", urlGetRepoPRs)
	assert.EqualValues(t, "https://api.github.com/repos/%s/%s/pulls/%s", urlGetRepoSinglePR)

}

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestGetRepoPRsErrorRestclient(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)

}

func TestGetRepoPRsErrorInvalidResponseBody(t *testing.T) {
	restclient.FlushMockups()
	invalidCloser, _ := os.Open("-asf3")
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidCloser,
		},
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body", err.Message)

}

func TestGetRepoPRsErrorInvalidErrorInterface(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body", err.Message)

}

func TestGetRepoPRsUnauthorized(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)

}

func TestGetRepoPRsInvalidResponse(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)

}

func TestGetRepoPRsNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})
	response, err := GetRepoPRs("", "test", "user1", "all")
	fmt.Println("some test")
	fmt.Println(response)
	fmt.Println(reflect.TypeOf(response).String())
	createDate, _ := time.Parse(time.RFC3339, "2019-11-27T14:30:10.578255Z")
	updateDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	closeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	mergeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")

	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.EqualValues(t, len(response), 1)
	assert.EqualValues(t, 123456, response[0].ID)
	assert.EqualValues(t, 9, response[0].Number)
	assert.EqualValues(t, "open", response[0].State)
	assert.EqualValues(t, "Title of the PR", response[0].Title)
	assert.EqualValues(t, createDate, response[0].CreatedAt)
	assert.EqualValues(t, updateDate, response[0].UpdatedAt)
	assert.EqualValues(t, closeDate, response[0].ClosedAt)
	assert.EqualValues(t, mergeDate, response[0].MergedAt)
	assert.EqualValues(t, "ABCDEF1234567890", response[0].MergeCommitSHA)

	assert.EqualValues(t, "My Login ID", response[0].User.Login)
	assert.EqualValues(t, 123456, response[0].User.ID)
	assert.EqualValues(t, "A user", response[0].User.Type)
	assert.EqualValues(t, true, response[0].User.SiteAdmin)

	assert.EqualValues(t, "A Second Login ID", response[0].Assignee.Login)
	assert.EqualValues(t, 8767, response[0].Assignee.ID)
	assert.EqualValues(t, "A user", response[0].Assignee.Type)
	assert.EqualValues(t, false, response[0].Assignee.SiteAdmin)

	assert.EqualValues(t, "A label", response[0].Base.Label)
	assert.EqualValues(t, "A Reference", response[0].Base.Ref)
	assert.EqualValues(t, "ABCDEF123456768", response[0].Base.SHA)

}

func TestGetRepoSinglePRErrorRestclient(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)

}

func TestGetRepoSinglePRErrorInvalidResponseBody(t *testing.T) {
	restclient.FlushMockups()
	invalidCloser, _ := os.Open("-asf3")
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidCloser,
		},
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body", err.Message)

}

func TestGetRepoSinglePRErrorInvalidErrorInterface(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body", err.Message)

}

func TestGetRepoSinglePRUnauthorized(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)

}

func TestGetRepoSinglePRInvalidResponse(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)

}

func TestGetRepoSinglePRNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`)),
		},
	})
	response, err := GetRepoSinglePR("", "test", "user1", "1")
	createDate, _ := time.Parse(time.RFC3339, "2019-11-27T14:30:10.578255Z")
	updateDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	closeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	mergeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", response.URL)
	assert.EqualValues(t, 123456, response.ID)
	assert.EqualValues(t, 9, response.Number)
	assert.EqualValues(t, "open", response.State)
	assert.EqualValues(t, "Title of the PR", response.Title)
	assert.EqualValues(t, createDate, response.CreatedAt)
	assert.EqualValues(t, updateDate, response.UpdatedAt)
	assert.EqualValues(t, closeDate, response.ClosedAt)
	assert.EqualValues(t, mergeDate, response.MergedAt)
	assert.EqualValues(t, "ABCDEF1234567890", response.MergeCommitSHA)

	assert.EqualValues(t, "My Login ID", response.User.Login)
	assert.EqualValues(t, 123456, response.User.ID)
	assert.EqualValues(t, "A user", response.User.Type)
	assert.EqualValues(t, true, response.User.SiteAdmin)

	assert.EqualValues(t, "A Second Login ID", response.Assignee.Login)
	assert.EqualValues(t, 8767, response.Assignee.ID)
	assert.EqualValues(t, "A user", response.Assignee.Type)
	assert.EqualValues(t, false, response.Assignee.SiteAdmin)

	assert.EqualValues(t, "A label", response.Base.Label)
	assert.EqualValues(t, "A Reference", response.Base.Ref)
	assert.EqualValues(t, "ABCDEF123456768", response.Base.SHA)

}

func TestGetGetSingleCommitPRErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetSingleCommitPR("", "test", "user1", "sha123")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)
}

func TestGetGetSingleCommitPRErrorResponseBody(t *testing.T) {

	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})
	response, err := GetSingleCommitPR("", "test", "user1", "sha123")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)
}

func TestGetGetSingleCommitPRNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})
	response, err := GetSingleCommitPR("", "test", "user1", "sha123")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, response[0].URL, "some URL")
	assert.EqualValues(t, response[0].ID, 123456)
	//the JSON is tested elswhere so not doing a full set of assertions here
}
