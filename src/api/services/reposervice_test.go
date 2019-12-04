package services

import (
	"io/ioutil"
	"net/http"
	"os"
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

func TestGetPRsInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("", "valid", "open")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid owner parameter", err.Message())

}

func TestGetPRsInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "", "open")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repo parameter", err.Message())
}

func TestGetPRsInvalidEmptyState(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid scope parameter", err.Message())
}

func TestGetPRsInvalidStateValue(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "repo", "some")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid scope parameter", err.Message())
}

func TestValidateInpusValidStateValues(t *testing.T) {
	owner, scope, state, err := validateInputs("owner", "repo", "open")
	assert.NotNil(t, owner)
	assert.NotNil(t, scope)
	assert.NotNil(t, state)
	assert.Nil(t, err)

	owner, scope, state, err = validateInputs("owner", "repo", "closed")
	assert.NotNil(t, owner)
	assert.NotNil(t, scope)
	assert.NotNil(t, state)
	assert.Nil(t, err)

	owner, scope, state, err = validateInputs("owner", "repo", "all")
	assert.NotNil(t, owner)
	assert.NotNil(t, scope)
	assert.NotNil(t, state)
	assert.Nil(t, err)
}

func TestGetPRsErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})

	response, err := RepositoryService.GetRepoPRs("test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())

}

func TestGetPRsNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})
	response, err := RepositoryService.GetRepoPRs("test", "user1", "all")
	createDate, _ := time.Parse(time.RFC3339, "2019-11-27T14:30:10.578255Z")
	updateDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	closeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	mergeDate, _ := time.Parse(time.RFC3339, "2019-10-28T14:30:10.578369Z")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", response[0].URL)
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

func TestRepoSinglePRInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("", "valid", "1")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid owner parameter", err.Message())

}

func TestRepoSinglePRInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "", "1")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repo parameter", err.Message())

}

func TestRepoSinglePRInvalidPR(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid pull parameter", err.Message())

}

func TestRepoSinglePRNotNumberPR(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid pull parameter", err.Message())

}

func TestRepoSinglePRErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})

	response, err := RepositoryService.GetRepoSinglePR("test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())

}

func TestRepoSinglePRNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`)),
		},
	})
	response, err := RepositoryService.GetRepoSinglePR("test", "user1", "1")
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
