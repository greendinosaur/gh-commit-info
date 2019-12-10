package githubprovider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greendinosaur/gh-commit-info/src/api/clients/restclient"
	"github.com/greendinosaur/gh-commit-info/src/api/domain/github"
)

//common constants and functions needed to access the Github API
const (
	//authorization header data, common to all github requests
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"

	//specific API header information, different Github API end points require different headers
	headerAccept              = "Accept"
	headerPRDraftAPI          = "application/vnd.github.shadow-cat-preview+json"
	headerPRForCommitDraftAPI = "application/vnd.github.groot-preview+json"
)

func getAuthorizationHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func getCommonHeader(accessToken string) http.Header {
	//All of the Github API calls require the token in the header
	//some may require additional header data which can be added
	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(accessToken))
	return headers
}

//getDataFromGitHub calls the Github API as indicated by the URL with the provided headers
//the headers will vary depending on the API end-point
//if there is an error in calling and parsing the response then an error is returned
//otherwise the response body is converted into bytes which can be unmarshalled into the relevant
//struct by the calling function
func getDataFromGithubAPI(URL string, headers http.Header) ([]byte, *github.GithubErrorResponse) {
	//the basic approach to calling Github to retrieve commit and PR data is the same irrespective
	//of the Github API being called and the data being returned
	//as a result, have put this logic into a common function
	response, err := restclient.Get(URL, headers)
	if err != nil {
		log.Println(fmt.Sprintf("error when calling Github API: %s", err.Error()))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid response body"}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid json response body"}
		}

		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	//all good so can return the body to be unmarshalled
	return bytes, nil
}

//getUnmarshalBodyError returns an error indicating there was a problem unmarhsalling the github response
func getUnmarshalBodyError() *github.GithubErrorResponse {
	return &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError,
		Message: "error when trying to unmarshal github response"}
}
