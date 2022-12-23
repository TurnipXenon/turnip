package models

import "net/http"

type ErrorWrapper struct {
	Err                 error `binding:"required"`
	UserMessage         string
	ShouldDisplayToUser bool
	HttpErrorCode       int
}

func (err ErrorWrapper) Error() string {
	return err.Err.Error()
}

func (err ErrorWrapper) Unwrap() error {
	return err.Err
}

func (err *ErrorWrapper) WriteHttpError(w http.ResponseWriter) {
	http.Error(w, err.UserMessage, err.HttpErrorCode)
}
