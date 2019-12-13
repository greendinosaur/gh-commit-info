package errors

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	myapiError := apiError{
		AStatus:  http.StatusInternalServerError,
		AMessage: "Some message",
		AnError:  "Some error",
	}

	bytes, err := json.Marshal(myapiError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.EqualValues(t, myapiError.AStatus, myapiError.Status())
	assert.EqualValues(t, myapiError.AMessage, myapiError.Message())
	assert.EqualValues(t, myapiError.AnError, myapiError.Error())

	var target apiError

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.NotNil(t, &target)
	assert.EqualValues(t, myapiError.AStatus, target.AStatus)
	assert.EqualValues(t, myapiError.AMessage, target.AMessage)
	assert.EqualValues(t, myapiError.AnError, target.AnError)
}

func TestNewNotFoundAPIError(t *testing.T) {
	apiError := NewNotFoundAPIError("Some Message")
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "Some Message", apiError.Message())
	assert.EqualValues(t, http.StatusNotFound, apiError.Status())
}

func TestNewInternalServerError(t *testing.T) {
	apiError := NewInternalServerError("Some Message")
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "Some Message", apiError.Message())
	assert.EqualValues(t, http.StatusInternalServerError, apiError.Status())
}

func TestNewBadRequestError(t *testing.T) {
	apiError := NewBadRequestError("Some Message")
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "Some Message", apiError.Message())
	assert.EqualValues(t, http.StatusBadRequest, apiError.Status())
}

func TestNewAPIErrorFromBytesValidJSON(t *testing.T) {
	myapiError := apiError{
		AStatus:  http.StatusInternalServerError,
		AMessage: "Some message",
		AnError:  "Some error",
	}

	bytes, err := json.Marshal(myapiError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	apiError, err := NewAPIErrorFromBytes(bytes)
	assert.Nil(t, err)
	assert.NotNil(t, apiError)
	assert.EqualValues(t, myapiError.Status(), apiError.Status())
	assert.EqualValues(t, myapiError.Message(), apiError.Message())
	assert.EqualValues(t, myapiError.Error(), apiError.Error())
}

func TestNewAPIErrorFromBytesInValidJSON(t *testing.T) {

	bytes, err := json.Marshal("{test:123,}")
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	apiError, err := NewAPIErrorFromBytes(bytes)
	assert.NotNil(t, err)
	assert.Nil(t, apiError)

}

func TestNewAPIError(t *testing.T) {
	apiError := NewAPIError(123, "Some Message")
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "Some Message", apiError.Message())
	assert.EqualValues(t, 123, apiError.Status())
}
