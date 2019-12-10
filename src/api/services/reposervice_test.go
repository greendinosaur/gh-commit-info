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
	owner, scope, state, err := validatePRInputs("owner", "repo", "open")
	assert.NotNil(t, owner)
	assert.NotNil(t, scope)
	assert.NotNil(t, state)
	assert.Nil(t, err)

	owner, scope, state, err = validatePRInputs("owner", "repo", "closed")
	assert.NotNil(t, owner)
	assert.NotNil(t, scope)
	assert.NotNil(t, state)
	assert.Nil(t, err)

	owner, scope, state, err = validatePRInputs("owner", "repo", "all")
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
	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
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
	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication

}

func TestSingleCommitPRNInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid owner parameter", err.Message())

}

func TestSingleCommitPRInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("owner", "", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repo parameter", err.Message())

}

func TestSingleCommitPRInvalidSHA(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid SHA parameter", err.Message())

}

func TestSingleCommitPRErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/ABC/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "ABC")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestSingleCommitPRNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "sha123")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, response[0].URL, "some URL")
	assert.EqualValues(t, response[0].ID, 123456)
	//the JSON is tested elswhere so not doing a full set of assertions here
}

func TestGetRepoCommitsInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoCommits("", "repo")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid owner parameter", err.Message())

}

func TestGetRepoCommitsInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoCommits("owner", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repo parameter", err.Message())

}

func TestGetRepoCommitsErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})

	response, err := RepositoryService.GetRepoCommits("test", "user1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestGetRepoCommitsNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`)),
		},
	})

	response, err := RepositoryService.GetRepoCommits("test", "user1")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, len(response), 1)
	assert.EqualValues(t, response[0].URL, "http://www.github.com")
	assert.EqualValues(t, response[0].SHA, "AABCDEF123456")

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, response[0].Commit.URL, "http://www.github.com")
}

func TestRepoCommitInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid owner parameter", err.Message())

}

func TestRepoSingleCommitInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("owner", "", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repo parameter", err.Message())

}

func TestRepoSingleCommitInvalidSHA(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid SHA parameter", err.Message())

}

func TestRepoSingleCommitErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestRepoSingleCommitNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}`)),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")

	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, response.URL, "http://www.github.com")
	assert.EqualValues(t, response.SHA, "AABCDEF123456")

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, response.Commit.URL, "http://www.github.com")
}
