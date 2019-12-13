package testutils

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

//GetMockDataUnauthorisedResponseStatusCode represents a mocked status code for unauthorised access
func GetMockDataUnauthorisedResponseStatusCode() int {
	return http.StatusUnauthorized
}

//GetMockDataUnauthorisedResponseMessage represents mock data for a github call that requires authentication
func GetMockDataUnauthorisedResponseMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`))
}

//GetMockDataPRsResponseStatusCode represents mock data to be used for a successful status code
func GetMockDataPRsResponseStatusCode() int {
	return http.StatusOK
}

//GetMockDataPRsResponseMessage represents mock data to be used for multiple PRs
func GetMockDataPRsResponseMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`))

}

//GetMockDataSinglePRResponseStatusCode represents mock data to be used for a successful status code
func GetMockDataSinglePRResponseStatusCode() int {
	return http.StatusOK
}

//GetMockDataSinglePRResponseMessage represents mock data to be used fo a single PR returned from a function
func GetMockDataSinglePRResponseMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`))
}

//GetMockDataCommitsResponseStatusCode represents mock data to be used for a successful status code
func GetMockDataCommitsResponseStatusCode() int {
	return http.StatusOK
}

//GetMockDataCommitsResponseMessage represents mock data to be used for multiple Commits returned from a function call
func GetMockDataCommitsResponseMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`))
}

//GetMockDataSingleCommitResponseStatusCode represents mock data to be used for a status code of a message returning a SingleCommitResponse struct
func GetMockDataSingleCommitResponseStatusCode() int {
	return http.StatusOK
}

//GetMockDataSingleCommitResponseMessage represents mock data to be used for a SingleCommitResponse
func GetMockDataSingleCommitResponseMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}`))
}

//GetMockDataSingleSliceCommitResponsesMessage returns a single commit that is a merge commit
func GetMockDataSingleSliceCommitResponsesMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`))
}

//GetMockDataSingleSliceNonMergeCommitResponsesMessage returns a commit that is not a merge commit
func GetMockDataSingleSliceNonMergeCommitResponsesMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"}]}]`))
}

//GetMockDataNoPRForCommitResponsesMessage returns no matching PRs
func GetMockDataNoPRForCommitResponsesMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[]`))
}

//GetMockDataApprovedPRForCommitResponsesMessage returns a matching approved PR for the commit sha AABCDEF123456
func GetMockDataApprovedPRForCommitResponsesMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`[{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"AABCDEF123456","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`))
}

//GetMockDataClosedPRForCommitResponsesMessage returns a matching closed but not merged PR for the commit sha AABCDEF123456
func GetMockDataClosedPRForCommitResponsesMessage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`))
}

//represents the error messages returned when validating parameters provided to service functions
const (
	ErrorMessageAuthentication = "Requires authentication"
	ErrorMessageRepo           = "invalid repo parameter"
	ErrorMessageOwner          = "invalid owner parameter"
	ErrorMessageScope          = "invalid scope parameter"
	ErrorMessagePull           = "invalid pull parameter"
	ErrorMessageInvalidSHA     = "invalid SHA parameter"
)
