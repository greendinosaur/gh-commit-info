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
	//r.GET("/repos/:owner/:repo/pulls", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"hello": "world"}) })
	c.Request = request
	//c.Params = make([]gin.Param, len(params))

	for paramName, paramValue := range params {
		c.Params = append(c.Params, gin.Param{Key: paramName, Value: paramValue})
	}
	//c.Params[0] = gin.Param{
	//	Key:   "owner",
	//	Value: "myowner",
	//}
	//c.Params[1] = gin.Param{
	//	Key:   "repo",
	//	Value: "myrepo",
	//}
	//c.Params[2] = gin.Param{
	//	Key:   "state",
	//	Value: "All",
	//}

	return c, r

}
