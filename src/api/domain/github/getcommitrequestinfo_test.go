package github

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParent(t *testing.T) {
	parent := Parent{
		URL: "http://test.com",
		SHA: "ABCDEF123456768",
	}

	bytes, err := json.Marshal(parent)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target Parent

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, parent.URL, target.URL)
	assert.EqualValues(t, parent.SHA, target.SHA)
}

func TestCommitUser(t *testing.T) {
	commitUser := CommitUser{
		Name:  "some name",
		Email: "email@email.com",
		Date:  time.Now().UTC(),
	}

	bytes, err := json.Marshal(commitUser)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target CommitUser

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, commitUser.Name, target.Name)
	assert.EqualValues(t, commitUser.Email, target.Email)
	assert.EqualValues(t, commitUser.Date, target.Date)
}

func TestDetailedCommit(t *testing.T) {

	author := CommitUser{
		Name:  "some name",
		Email: "email@email.com",
		Date:  time.Now().UTC(),
	}

	committer := CommitUser{
		Name:  "some committer",
		Email: "someemail@email.com",
		Date:  time.Now().UTC(),
	}

	detailedCommit := DetailedCommitInfo{
		URL:       "http://www.github.com",
		Author:    author,
		Committer: committer,
		Message:   "some commit message",
	}

	bytes, err := json.Marshal(detailedCommit)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target DetailedCommitInfo

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, detailedCommit.URL, target.URL)
	assert.EqualValues(t, detailedCommit.Message, detailedCommit.Message)
	assert.EqualValues(t, detailedCommit.Author, target.Author)
	assert.EqualValues(t, detailedCommit.Committer, target.Committer)

}

func TestGetCommitInfo(t *testing.T) {
	authorGitUser := GitUser{
		Login:     "some loing id",
		ID:        9876,
		Type:      "user",
		SiteAdmin: true,
	}

	committerGitUser := GitUser{
		Login:     "login id",
		ID:        12345,
		Type:      "user",
		SiteAdmin: false,
	}

	author := CommitUser{
		Name:  "some name",
		Email: "email@email.com",
		Date:  time.Now().UTC(),
	}

	committer := CommitUser{
		Name:  "some committer",
		Email: "someemail@email.com",
		Date:  time.Now().UTC(),
	}

	detailedCommit := DetailedCommitInfo{
		URL:       "http://www.github.com",
		Author:    author,
		Committer: committer,
		Message:   "some commit message",
	}

	parent1 := Parent{
		URL: "http://test.com",
		SHA: "ABCDEF123456768",
	}

	parent2 := Parent{
		URL: "http://test12.com",
		SHA: "ABFGGG",
	}

	var parents []Parent
	parents = append(parents, parent1, parent2)

	getCommitInfo := GetCommitInfo{
		URL:       "http://www.github.com",
		SHA:       "AABCDEF123456",
		Commit:    detailedCommit,
		Author:    authorGitUser,
		Committer: committerGitUser,
		Parents:   parents,
	}

	bytes, err := json.Marshal(getCommitInfo)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target GetCommitInfo

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, getCommitInfo.URL, target.URL)
	assert.EqualValues(t, getCommitInfo.SHA, target.SHA)
	assert.EqualValues(t, getCommitInfo.Commit, target.Commit)
	assert.EqualValues(t, getCommitInfo.Author, target.Author)
	assert.EqualValues(t, getCommitInfo.Committer, target.Committer)
	assert.EqualValues(t, getCommitInfo.Parents, target.Parents)

}
