package bobby

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//need to fix this import
	"github.com/greendinosaur/golang-microservices/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "bobby and his chariots", chariot)
}

func TestBobby(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bobby", nil)
	c := test_utils.GetMockedContext(request, response)

	Chariot(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "bobby and his chariots", response.Body.String())
}
