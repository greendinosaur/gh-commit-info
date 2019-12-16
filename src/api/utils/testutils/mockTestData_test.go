package testutils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, ErrorMessageAuthentication, "Requires authentication")
	assert.EqualValues(t, ErrorMessageRepo, "invalid repo parameter")
	assert.EqualValues(t, ErrorMessageAuthentication, "Requires authentication")
	assert.EqualValues(t, ErrorMessageOwner, "invalid owner parameter")
	assert.EqualValues(t, ErrorMessageScope, "invalid scope parameter")
	assert.EqualValues(t, ErrorMessagePull, "invalid pull parameter")
	assert.EqualValues(t, ErrorMessageInvalidSHA, "invalid SHA parameter")

}

func TestGetMockDataUnauthorisedResponseStatusCode(t *testing.T) {
	assert.EqualValues(t, http.StatusUnauthorized, GetMockDataUnauthorisedResponseStatusCode())
}

func TestGetMockDataPRsResponseStatusCode(t *testing.T) {
	assert.EqualValues(t, http.StatusOK, GetMockDataPRsResponseStatusCode())
}

func TestGetMockDataSinglePRResponseStatusCode(t *testing.T) {
	assert.EqualValues(t, http.StatusOK, GetMockDataSinglePRResponseStatusCode())
}

func TestGetMockDataCommitsResponseStatusCode(t *testing.T) {
	assert.EqualValues(t, http.StatusOK, GetMockDataCommitsResponseStatusCode())
}

func TestGetMockDataSingleCommitResponseStatusCode(t *testing.T) {
	assert.EqualValues(t, http.StatusOK, GetMockDataSingleCommitResponseStatusCode())
}

func TestGetMockDataUnauthorisedResponseMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataUnauthorisedResponseMessage())
	newStr := buf.String()
	assert.EqualValues(t, `{"message": "Requires authentication"}`, newStr)
}

func TestGetMockDataPRsResponseMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataPRsResponseMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`, newStr)
}

func TestGetMockDataSinglePRResponseMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataSinglePRResponseMessage())
	newStr := buf.String()
	assert.EqualValues(t, `{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}`, newStr)
}

func TestGetMockDataCommitsResponseMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataCommitsResponseMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`, newStr)
}

func TestGetMockDataSingleCommitResponseMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataSingleCommitResponseMessage())
	newStr := buf.String()
	assert.EqualValues(t, `{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}`, newStr)
}

func TestGetMockDataSingleSliceCommitResponsesMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataSingleSliceCommitResponsesMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"},{"url":"http://test12.com","sha":"ABFGGG"}]}]`, newStr)
}

func TestGetMockDataSingleSliceNonMergeCommitResponsesMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataSingleSliceNonMergeCommitResponsesMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"http://www.github.com","sha":"AABCDEF123456","commit":{"url":"http://www.github.com","author":{"name":"some name","email":"email@email.com","date":"2019-12-09T15:00:04.061358Z"},"committer":{"name":"some committer","email":"someemail@email.com","date":"2019-12-09T15:00:04.061358Z"},"message":"some commit message"},"author":{"login":"some loing id","id":9876,"type":"user","site_admin":true},"committer":{"login":"login id","id":12345,"type":"user","site_admin":false},"parents":[{"url":"http://test.com","sha":"ABCDEF123456768"}]}]`, newStr)
}

func TestGetMockDataNoPRForCommitResponsesMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataNoPRForCommitResponsesMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[]`, newStr)
}

func TestGetMockDataApprovedPRForCommitResponsesMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataApprovedPRForCommitResponsesMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"AABCDEF123456","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`, newStr)
}

func TestGetMockDataClosedPRForCommitResponsesMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(GetMockDataClosedPRForCommitResponsesMessage())
	newStr := buf.String()
	assert.EqualValues(t, `[{"url":"some URL","id":123456,"number":9,"state":"closed","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`, newStr)
}
