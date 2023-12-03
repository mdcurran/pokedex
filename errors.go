package pokedex

import (
	"fmt"
	"net/http"
)

const (
	// CodeInternal is an error code that should be returned when an error
	// occurs that's not in the HTTP request lifecycle.
	CodeInternal int = 600
	// CodeClientClosed is returned when client.Close() has been previously
	// called. A client that is closed is unable to cache API responses, so
	// we want to force users to create a new client.
	CodeClientClosed int = 601
	// CodeInvalidArgs indicate an API request using the SDK has malformed
	// arguments.
	CodeInvalidArgs int = 602
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
