package githubdomain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubError(t *testing.T) {
	ghError := GithubError{
		Resource: "some resource",
		Code:     "code123",
		Field:    "some field",
		Message:  "an error message",
	}

	bytes, err := json.Marshal(ghError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target GithubError

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, ghError.Resource, target.Resource)
	assert.EqualValues(t, ghError.Code, target.Code)
	assert.EqualValues(t, ghError.Field, target.Field)
	assert.EqualValues(t, ghError.Message, target.Message)

}

func TestGithubErrorResponse(t *testing.T) {
	ghError := GithubError{
		Resource: "some resource",
		Code:     "code123",
		Field:    "some field",
		Message:  "an error message",
	}

	ghErrorResponse := GithubErrorResponse{
		StatusCode:       501,
		Message:          "an error message",
		DocumentationURL: "www.github.com",
		Errors:           []GithubError{ghError},
	}

	bytes, err := json.Marshal(ghErrorResponse)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target GithubErrorResponse
	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, ghErrorResponse.StatusCode, target.StatusCode)
	assert.EqualValues(t, ghErrorResponse.Message, target.Message)
	assert.EqualValues(t, ghErrorResponse.DocumentationURL, target.DocumentationURL)
	assert.EqualValues(t, ghErrorResponse.Errors[0], target.Errors[0])

}

func TestError(t *testing.T) {
	ghErrorResponse := GithubErrorResponse{
		StatusCode:       501,
		Message:          "an error message",
		DocumentationURL: "www.github.com",
	}

	assert.EqualValues(t, ghErrorResponse.Message, ghErrorResponse.Error())
}
