package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

const (
	missingMock = "no mockup found for given request"
)

//Mock is used to mock up responses given a request
type Mock struct {
	URL        string
	HTTPMethod string
	Response   *http.Response
	Err        error
}

func getMockID(HTTPMethod string, URL string) string {
	return fmt.Sprintf("%s_%s", HTTPMethod, URL)
}

//StartMockups indicates mocks are to be used
func StartMockups() {
	enabledMocks = true
}

//StopMockups stops mocks from being used
func StopMockups() {
	enabledMocks = false
}

//FlushMockups resets the mocks to empty
func FlushMockups() {
	mocks = make(map[string]*Mock)
}

//AddMockup adds a new mockup
func AddMockup(mock Mock) {
	mocks[getMockID(mock.HTTPMethod, mock.URL)] = &mock
}

func getMockResponse(method string, URL string) (*http.Response, error) {
	mock := mocks[getMockID(method, URL)]

	if mock == nil {
		return nil, errors.New(missingMock)
	}

	return mock.Response, mock.Err
}

func getContent(method string, URL string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, URL, bytes.NewReader(jsonBytes))
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}

//Post submits a HTTP POST request with the given parameters
func Post(URL string, body interface{}, headers http.Header) (*http.Response, error) {

	if enabledMocks {
		return getMockResponse(http.MethodPost, URL)
	}

	return getContent(http.MethodPost, URL, body, headers)
}

//Get submits a HTTP GET request with the given parameters
func Get(URL string, headers http.Header) (*http.Response, error) {

	if enabledMocks {
		return getMockResponse(http.MethodGet, URL)
	}

	return getContent(http.MethodGet, URL, nil, headers)

}
