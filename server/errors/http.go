package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Msg        string
	StatusCode int
	Err        error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func NewHTTPError(msg string, statusCode int, err error) *HTTPError {
	return &HTTPError{
		Msg:        msg,
		StatusCode: statusCode,
		Err:        err,
	}
}

/*
Convenience method meant to check whether or not err is an instance of HTTPError and to set the response status accordingly.

This function is meant to replace the usage of net/http's http.Error function.
*/
func APIError(w http.ResponseWriter, err error) {
	var httpError *HTTPError
	if errors.As(err, &httpError) {
		http.Error(w, err.Error(), httpError.StatusCode)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
