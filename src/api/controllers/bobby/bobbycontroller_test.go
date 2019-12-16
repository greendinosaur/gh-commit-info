package bobby

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greendinosaur/gh-commit-info/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "bobby and his chariots", chariot)
}

func TestBobby(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bobby", nil)
	c, _ := testutils.GetMockedContext(request, response)

	Chariot(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "bobby and his chariots", response.Body.String())
}
