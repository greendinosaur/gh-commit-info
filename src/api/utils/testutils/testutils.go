package testutils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

//GetMockedContext returns a test context to be used in tests
func GetMockedContext(request *http.Request, response *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	c, r := gin.CreateTestContext(response)
	c.Request = request
	return c, r

}

//GetMockedContextWithParams returns a test context to be used in tests with Params set to
//the key/value pairs in the provided map
func GetMockedContextWithParams(request *http.Request, response *httptest.ResponseRecorder, params map[string]string) (*gin.Context, *gin.Engine) {
	c, r := gin.CreateTestContext(response)
	c.Request = request

	for paramName, paramValue := range params {
		c.Params = append(c.Params, gin.Param{Key: paramName, Value: paramValue})
	}
	return c, r

}
