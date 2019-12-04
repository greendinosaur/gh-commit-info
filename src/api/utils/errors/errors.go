package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

//APIError defines an error in using the API
type APIError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	AStatus  int    `json:"status"`
	AMessage string `json:"message"`
	AnError  string `json:"error,omitempty"`
}

//Status returns the status code
func (e *apiError) Status() int {
	return e.AStatus
}

//Message returns the error message
func (e *apiError) Message() string {
	return e.AMessage
}

//Error returns the error
func (e *apiError) Error() string {
	return e.AnError
}

//NewNotFoundAPIError returns an APIError indicating data wasn't found
func NewNotFoundAPIError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusNotFound,
		AMessage: message,
	}
}

//NewInternalServerError retuns an error indicating an internal server error
func NewInternalServerError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusInternalServerError,
		AMessage: message,
	}
}

//NewBadRequestError returns an error indicating a problem with the API request
func NewBadRequestError(message string) APIError {
	return &apiError{
		AStatus:  http.StatusBadRequest,
		AMessage: message,
	}
}

//NewAPIErrorFromBytes returns an error based on a byte slice
func NewAPIErrorFromBytes(body []byte) (APIError, error) {
	var result apiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("invalid json body")
	}

	return &result, nil

}

//NewAPIError returns an error with the provided inputs
func NewAPIError(statusCode int, message string) APIError {
	return &apiError{
		AStatus:  statusCode,
		AMessage: message,
	}
}
