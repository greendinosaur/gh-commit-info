package githubprovider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/stretchr/testify/assert"
)

func TestConstantsForCommits(t *testing.T) {
	assert.EqualValues(t, "https://api.github.com/repos/%s/%s/commits", urlGetRepoCommits)
	assert.EqualValues(t, "https://api.github.com/repos/%s/%s/commits/%s", urlGetRepoSingleCommit)
}

func TestGetRepoCommitsErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits",
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetRepoCommits("", "myuser", "myrepo")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)
}

func TestGetRepoCommitsErrorResponseBody(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})
	response, err := GetRepoCommits("", "myuser", "myrepo")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)
}

func TestGetRepoCommitsNoError(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`)),
		},
	})
	response, err := GetRepoCommits("", "myuser", "myrepo")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, len(response), 1)
	assert.EqualValues(t, response[0].URL, "http://www.github.com")
	assert.EqualValues(t, response[0].SHA, "AABCDEF123456")

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, "http://www.github.com", response[0].Commit.URL)

}

func TestGetRepoCommitsDateRangeErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(FmtGithubDate) + "&until=" + toDate.UTC().Format(FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetRepoCommitsInDateRange("", "myuser", "myrepo", fromDate, toDate)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)
}

func TestGetRepoCommitsDateRangeErrorResponseBody(t *testing.T) {

	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(FmtGithubDate) + "&until=" + toDate.UTC().Format(FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})
	response, err := GetRepoCommitsInDateRange("", "myuser", "myrepo", fromDate, toDate)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)
}

func TestGetRepoCommitsDateRangeNoError(t *testing.T) {

	restclient.FlushMockups()
	fromDate := time.Now().UTC().AddDate(-1, 0, 0)
	toDate := time.Now().UTC()
	urlForMock := "https://api.github.com/repos/myuser/myrepo/commits?since=" + fromDate.UTC().Format(FmtGithubDate) + "&until=" + toDate.UTC().Format(FmtGithubDate)

	restclient.AddMockup(restclient.Mock{
		URL:        urlForMock,
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`)),
		},
	})

	response, err := GetRepoCommitsInDateRange("", "myuser", "myrepo", fromDate, toDate)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, len(response), 1)
	assert.EqualValues(t, response[0].URL, "http://www.github.com")
	assert.EqualValues(t, response[0].SHA, "AABCDEF123456")

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, "http://www.github.com", response[0].Commit.URL)

}

func TestGetRepoSingleCommitErrorFromGithub(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits/abcdef123",
		HTTPMethod: http.MethodGet,
		Err:        errors.New("invalid rest client response"),
	})
	response, err := GetRepoSingleCommit("", "myuser", "myrepo", "abcdef123")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response", err.Message)

}

func TestGetRepoSingleCommitErrorResponseBody(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits/abcdef123",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`[{"id": "123"}]`)),
		},
	})
	response, err := GetRepoSingleCommit("", "myuser", "myrepo", "abcdef123")
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error when trying to unmarshal github response", err.Message)
}

func TestGetRepoSingleCommitNoError(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL:        "https://api.github.com/repos/myuser/myrepo/commits/abcdef123",
		HTTPMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}`)),
		},
	})
	response, err := GetRepoSingleCommit("", "myuser", "myrepo", "abcdef123")
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "http://www.github.com", response.URL)
	assert.EqualValues(t, "AABCDEF123456", response.SHA)

	//will just test one of the items has been marshalled okay
	//the domain object tests this fully so no need for duplication
	assert.EqualValues(t, "http://www.github.com", response.Commit.URL)

}
