package bobby

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	chariot = "bobby and his chariots"
)

//Chariot sets the context with an OK status and a message
//can be used to check the web server is up OK
func Chariot(c *gin.Context) {
	c.String(http.StatusOK, chariot)
}
