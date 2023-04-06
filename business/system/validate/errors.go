package validate

import (
	"encoding/json"
	"errors"
)

var ErrInvalidID = errors.New("ID is not in proper form")

// ErrorReponse is the form used for API response from failures in API.
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

// RequestError wraps a provided error with an HTTP status code. This function should be used when
// handlers encounter expected errors.
type RequestError struct {
	Err    error
	Status int
	Fields error
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

func (err *RequestError) Error() string {
	return err.Err.Error()
}

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Cauaase iterates through all wrapped errors until it reaches the root error
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}
