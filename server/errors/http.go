package errors

import "fmt"

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
