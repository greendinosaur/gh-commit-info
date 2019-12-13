package githubdomain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
