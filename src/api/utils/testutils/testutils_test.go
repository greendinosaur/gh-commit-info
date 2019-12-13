package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMockedConstant(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/something", nil)
	assert.Nil(t, err)
	request.Header = http.Header{"X-Mock": {"true"}}
	c, _ := GetMockedContext(request, response)

	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "8080", c.Request.URL.Port())
	assert.EqualValues(t, "/something", c.Request.URL.Path)
	assert.EqualValues(t, "http", c.Request.URL.Scheme)
	assert.EqualValues(t, 1, len(c.Request.Header))
	assert.EqualValues(t, "true", c.GetHeader("x-mock"))
	assert.EqualValues(t, "true", c.GetHeader("X-mock"))

}

func TestGetMockedConstantParams(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/something", nil)
	assert.Nil(t, err)
	request.Header = http.Header{"X-Mock": {"true"}}

	params := map[string]string{"owner": "myowner", "repo": "myrepo", "state": "all"}

	c, _ := GetMockedContextWithParams(request, response, params)

	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "8080", c.Request.URL.Port())
	assert.EqualValues(t, "/something", c.Request.URL.Path)
	assert.EqualValues(t, "http", c.Request.URL.Scheme)
	assert.EqualValues(t, 1, len(c.Request.Header))
	assert.EqualValues(t, "true", c.GetHeader("x-mock"))
	assert.EqualValues(t, "true", c.GetHeader("X-mock"))
	assert.EqualValues(t, "myowner", c.Param("owner"))
	assert.EqualValues(t, "myrepo", c.Param("repo"))
	assert.EqualValues(t, "all", c.Param("state"))

}


