package github

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepoBase(t *testing.T) {
	repoBase := RepoBase{
		Label: "A label",
		Ref:   "A Reference",
		SHA:   "ABCDEF123456768",
	}

	bytes, err := json.Marshal(repoBase)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target RepoBase

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, repoBase.Label, target.Label)
	assert.EqualValues(t, repoBase.Ref, target.Ref)
	assert.EqualValues(t, repoBase.SHA, target.SHA)
}

func TestGitUser(t *testing.T) {
	gitUser := GitUser{
		Login:     "My Login ID",
		ID:        123456,
		Type:      "A user",
		SiteAdmin: true,
	}

	bytes, err := json.Marshal(gitUser)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target GitUser

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, gitUser.Login, target.Login)
	assert.EqualValues(t, gitUser.ID, target.ID)
	assert.EqualValues(t, gitUser.Type, target.Type)
	assert.EqualValues(t, gitUser.SiteAdmin, target.SiteAdmin)
}

func TestGetPullRequestInfoResponse(t *testing.T) {

	repoBase := RepoBase{
		Label: "A label",
		Ref:   "A Reference",
		SHA:   "ABCDEF123456768",
	}

	gitUser1 := GitUser{
		Login:     "My Login ID",
		ID:        123456,
		Type:      "A user",
		SiteAdmin: true,
	}

	gitUser2 := GitUser{
		Login:     "A Second Login ID",
		ID:        8767,
		Type:      "A user",
		SiteAdmin: false,
	}
	getPRInfoResponse := GetSinglePullRequestResponse{
		URL:            "some URL",
		ID:             123456,
		Number:         9,
		State:          "open",
		Title:          "Title of the PR",
		CreatedAt:      time.Now().AddDate(0, 0, -1).UTC(),
		UpdatedAt:      time.Now().AddDate(0, -1, 0).UTC(),
		ClosedAt:       time.Now().AddDate(0, -1, 0).UTC(),
		MergedAt:       time.Now().AddDate(0, -1, 0).UTC(),
		MergeCommitSHA: "ABCDEF1234567890",
		User:           gitUser1,
		Assignee:       gitUser2,
		Base:           repoBase,
	}

	bytes, err := json.Marshal(getPRInfoResponse)
	jsonVal := string(bytes)
	t.Log(jsonVal)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target GetSinglePullRequestResponse

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, getPRInfoResponse.URL, target.URL)
	assert.EqualValues(t, getPRInfoResponse.ID, target.ID)
	assert.EqualValues(t, getPRInfoResponse.Number, target.Number)
	assert.EqualValues(t, getPRInfoResponse.State, target.State)
	assert.EqualValues(t, getPRInfoResponse.Title, target.Title)
	assert.EqualValues(t, getPRInfoResponse.CreatedAt, target.CreatedAt)
	assert.EqualValues(t, getPRInfoResponse.UpdatedAt, target.UpdatedAt)
	assert.EqualValues(t, getPRInfoResponse.ClosedAt, target.ClosedAt)
	assert.EqualValues(t, getPRInfoResponse.MergedAt, target.MergedAt)
	assert.EqualValues(t, getPRInfoResponse.MergeCommitSHA, target.MergeCommitSHA)
	assert.EqualValues(t, getPRInfoResponse.User, target.User)
	assert.EqualValues(t, getPRInfoResponse.Assignee, target.Assignee)
	assert.EqualValues(t, getPRInfoResponse.Base, target.Base)

}

func TestMultiplePullRequestResponse(t *testing.T) {

	jsonAsString := `[{"url":"some URL","id":123456,"number":9,"state":"open","title":"Title of the PR","created_at":"2019-11-27T14:30:10.578255Z","updated_at":"2019-10-28T14:30:10.578369Z","closed_at":"2019-10-28T14:30:10.578369Z","merged_at":"2019-10-28T14:30:10.578369Z","merge_commit_sha":"ABCDEF1234567890","user":{"login":"My Login ID","id":123456,"type":"A user","site_admin":true},"assignee":{"login":"A Second Login ID","id":8767,"type":"A user","site_admin":false},"base":{"label":"A label","ref":"A Reference","sha":"ABCDEF123456768"}}]`

	bytes := []byte(jsonAsString)

	target := []MultiplePullRequestResponse{}

	err := json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, "some URL", target[0].URL)

}
