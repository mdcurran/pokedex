package pokedex

import (
	"fmt"
	"net/http"
)

// CodeInternal is an error code that should be returned when an error
// occurs that's not in the HTTP request lifecycle.
const (
	CodeInternal     int = 600
	CodeClientClosed int = 601
)

type SDKError struct {
	Message    string
	StatusCode int
	Response   *http.Response
}

func NewError(message string, code int, response *http.Response) *SDKError {
	return &SDKError{
		Message:    message,
		StatusCode: code,
		Response:   response,
	}
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("message: %q status: %d", e.Message, e.StatusCode)
}
