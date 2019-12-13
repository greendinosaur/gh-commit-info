package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/providers/githubprovider"
	"github.com/greendinosaur/gh-commit-info/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

//these test the logic for checking parameters
func TestGetPRsInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("", "valid", "open")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())

}

func TestGetPRsInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "", "open")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())
}

func TestGetPRsInvalidEmptyState(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageScope, err.Message())
}

func TestGetPRsInvalidStateValue(t *testing.T) {
	result, err := RepositoryService.GetRepoPRs("owner", "repo", "some")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageScope, err.Message())
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

//these test the logic for getting multiple PRs
func TestGetPRsErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoPRs("test", "user1", "all")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())

}

func TestGetPRsNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls?state=all",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataPRsResponseStatusCode(),
			Body:       testutils.GetMockDataPRsResponseMessage(),
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

//these test the logic for getting a single PR

func TestRepoSinglePRInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("", "valid", "1")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())

}

func TestRepoSinglePRInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "", "1")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())

}

func TestRepoSinglePRInvalidPR(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessagePull, err.Message())

}

func TestRepoSinglePRNotNumberPR(t *testing.T) {
	result, err := RepositoryService.GetRepoSinglePR("valid", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessagePull, err.Message())

}

func TestRepoSinglePRErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoSinglePR("test", "user1", "1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())

}

func TestRepoSinglePRNoError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/pulls/1",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSinglePRResponseStatusCode(),
			Body:       testutils.GetMockDataSinglePRResponseMessage(),
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

//these test the logic for getting the PRs associated with a single commit
func TestSingleCommitPRInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())

}

func TestSingleCommitPRInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("owner", "", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())

}

func TestSingleCommitPRInvalidSHA(t *testing.T) {
	result, err := RepositoryService.GetSingleCommitPR("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageInvalidSHA, err.Message())

}

func TestSingleCommitPRErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/ABC/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "ABC")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
}

func TestSingleCommitPRNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataPRsResponseStatusCode(),
			Body:       testutils.GetMockDataPRsResponseMessage(),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "sha123")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "some URL", response[0].URL)
	assert.EqualValues(t, 123456, response[0].ID)
	//the JSON is tested elswhere so not doing a full set of assertions here
}

//these test the logic for getting multiple repo commits
func TestGetRepoCommitsInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoCommits("", "repo")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())

}

func TestGetRepoCommitsInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoCommits("owner", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())

}

func TestGetRepoCommitsErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoCommits("test", "user1")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
}

func TestGetRepoCommitsNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataCommitsResponseStatusCode(),
			Body:       testutils.GetMockDataCommitsResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoCommits("test", "user1")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(response))
	assert.EqualValues(t, "http://www.github.com", response[0].URL)
	assert.EqualValues(t, "AABCDEF123456", response[0].SHA)

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, "http://www.github.com", response[0].Commit.URL)
}

//these test the logic for getting a single commit for a repo
func TestRepoCommitInvalidOwner(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("", "repo", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())

}

func TestRepoSingleCommitInvalidRepo(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("owner", "", "asd")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())

}

func TestRepoSingleCommitInvalidSHA(t *testing.T) {
	result, err := RepositoryService.GetRepoSingleCommit("owner", "repo", "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageInvalidSHA, err.Message())

}

func TestRepoSingleCommitErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
}

func TestRepoSingleCommitNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataSingleCommitResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")

	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "http://www.github.com", response.URL)
	assert.EqualValues(t, "AABCDEF123456", response.SHA)

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, response.Commit.URL, "http://www.github.com")
}

//these test the logic for determining whether the commit is a merge commit
func TestMergeCommitIsMerge(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataSingleCommitResponseMessage(),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")
	response.IsMergeCommit = isMergeCommit(response)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, true, response.IsMergeCommit)

}

func TestMergeCommitIsNotMerge(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/shaabcd",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"}]}`)),
		},
	})

	response, err := RepositoryService.GetRepoSingleCommit("test", "user1", "shaabcd")

	response.IsMergeCommit = isMergeCommit(response)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, false, response.IsMergeCommit)
}

//these test the logic for determining whether a PR resuled in a merge commit
func TestIsPRResultingInMergeSuccess(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "sha123")

	PRResultsInMerge := isPRResultingInMerge(&response[0])
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, true, PRResultsInMerge)
}

func TestIsPRResultingInMergeFailInvalidPRState(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "sha123")

	PRResultsInMerge := isPRResultingInMerge(&response[0])
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, false, PRResultsInMerge)
}

func TestIsPRResultingInMergeFailInvalidMergeSHA(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/test/user1/commits/sha123/pulls",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`)),
		},
	})

	response, err := RepositoryService.GetSingleCommitPR("test", "user1", "sha123")

	PRResultsInMerge := isPRResultingInMerge(&response[0])
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, false, PRResultsInMerge)
}

//these test GetRepoCommits

func TestGetRepoCommitsInDateRangeInvalidOwner(t *testing.T) {
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()

	response, err := getRepoCommitsInDateRange("", "owner", fromDate, toDate)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageOwner, err.Message())
}

func TestGetRepoCommitsInDateRangeInvalidRepo(t *testing.T) {

	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()

	response, err := getRepoCommitsInDateRange("myuser", "", fromDate, toDate)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageRepo, err.Message())

}

func TestGetRepoCommitsInDateRangeErrorWithGithub(t *testing.T) {
	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(githubprovider.FmtGithubDate) + "&until=" + toDate.UTC().Format(githubprovider.FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := getRepoCommitsInDateRange("myuser", "myrepo", fromDate, toDate)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
}

func TestGetRepoCommitsInDateRangeSuccess(t *testing.T) {

	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(githubprovider.FmtGithubDate) + "&until=" + toDate.UTC().Format(githubprovider.FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`)),
		},
	})

	response, err := getRepoCommitsInDateRange("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, len(response), 1)
	assert.EqualValues(t, "http://www.github.com", response[0].URL)
	assert.EqualValues(t, "AABCDEF123456", response[0].SHA)

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, "http://www.github.com", response[0].URL)

}

//these test the GetCodeReviewReport function
func TestGetCodeReviewReportErrorGettingCommits(t *testing.T) {
	//simular getting an error when getting hold of the initial commits
	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(githubprovider.FmtGithubDate) + "&until=" + toDate.UTC().Format(githubprovider.FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetCodeReviewReport("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
	assert.EqualValues(t, "", response)
}

func TestGetCodeReviewReportErrorGettingPR(t *testing.T) {
	//simulate getting an error when getting hold of a PR associated aith a commit
	//need to add two mocks, one for getting the commit and one for getting its PRs
	//the one for the PRs needs to return an error
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
			StatusCode: testutils.GetMockDataUnauthorisedResponseStatusCode(),
			Body:       testutils.GetMockDataUnauthorisedResponseMessage(),
		},
	})

	response, err := RepositoryService.GetCodeReviewReport("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.NotNil(t, err) //need to check the error message
	assert.EqualValues(t, "", response)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, testutils.ErrorMessageAuthentication, err.Message())
}

func TestGetCodeReviewReportSuccessMergeCommit(t *testing.T) {
	//need to have test data where there is a merge commit
	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.Format(githubprovider.FmtGithubDate) + "&until=" + toDate.Format(githubprovider.FmtGithubDate)

	fmt.Println(urlForMock)
	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: testutils.GetMockDataSingleCommitResponseStatusCode(),
			Body:       testutils.GetMockDataSingleSliceCommitResponsesMessage(),
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

	response, err := RepositoryService.GetCodeReviewReport("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "#Total Commits: 1, #Merged Commits: 1,  #Commits with PRs: 1, #Non-merge commits with no PR: 0", response)
}

func TestGetCodeReviewReportSuccessCommitWithPR(t *testing.T) {
	//need to have test data where there is a commit with a PR
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

	response, err := RepositoryService.GetCodeReviewReport("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "#Total Commits: 1, #Merged Commits: 0,  #Commits with PRs: 1, #Non-merge commits with no PR: 0", response)

}

func TestGetCodeReviewReportSuccessCommitWithNoPR(t *testing.T) {
	//need to have test data where there is a commit with no PR
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
			Body:       testutils.GetMockDataNoPRForCommitResponsesMessage(),
		},
	})

	response, err := RepositoryService.GetCodeReviewReport("myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "#Total Commits: 1, #Merged Commits: 0,  #Commits with PRs: 0, #Non-merge commits with no PR: 1", response)

}
