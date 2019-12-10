package restclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartMockups(t *testing.T) {
	StartMockups()
	assert.EqualValues(t, true, enabledMocks)
}

func TestStopMockups(t *testing.T) {
	StopMockups()
	assert.EqualValues(t, false, enabledMocks)
}

func TestFlushMockups(t *testing.T) {
	FlushMockups()
	assert.EqualValues(t, 0, len(mocks))
}

func TestAddMockupMissing(t *testing.T) {
	URL := "https://api.github.com/repos/myowner/myrepo/pulls?state=all"
	method := http.MethodGet
	AddMockup(Mock{
		URL:        URL,
		HTTPMethod: method,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
		Err: nil,
	})

	mock := mocks[getMockID(method, URL+"&somerandomtext")]
	assert.Nil(t, mock)

}

func TestAddMockupPresent(t *testing.T) {

	URL := "https://api.github.com/repos/myowner/myrepo/pulls?state=all"
	method := http.MethodGet
	AddMockup(Mock{
		URL:        URL,
		HTTPMethod: method,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
		Err: nil,
	})

	mock := mocks[getMockID(method, URL)]
	assert.NotNil(t, mock)
	assert.EqualValues(t, URL, mock.URL)
	assert.EqualValues(t, method, mock.HTTPMethod)
	assert.Nil(t, mock.Err)
	assert.NotNil(t, mock.Response)
}

func TestGetMockResponseNotPresent(t *testing.T) {
	URL := "https://api.github.com/repos/myowner/myrepo/pulls?state=all"
	method := http.MethodGet
	AddMockup(Mock{
		URL:        URL,
		HTTPMethod: method,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
		Err: nil,
	})
	mock, err := getMockResponse(method, URL+"sometext")
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("%+v", err), missingMock)
	assert.Nil(t, mock)
}

func TestGetMockResponseMockPresent(t *testing.T) {
	URL := "https://api.github.com/repos/myowner/myrepo/pulls?state=all"
	method := http.MethodGet
	AddMockup(Mock{
		URL:        URL,
		HTTPMethod: method,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
		},
		Err: nil,
	})
	mock, err := getMockResponse(method, URL)
	assert.NotNil(t, mock)
	assert.Nil(t, err)
	assert.NotNil(t, mock)
	assert.EqualValues(t, http.StatusUnauthorized, mock.StatusCode)
	assert.EqualValues(t, ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)), mock.Body)
}
