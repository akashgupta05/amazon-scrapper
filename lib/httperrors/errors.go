package httperrors

import (
	"errors"
	"net/http"
)

type HttpErrorInterface interface {
}

type HttpError struct {
	StatusCode int
	Error      error `json:"error"`
}

func NewHttpError(message string) *HttpError {
	return &HttpError{
		Error:      errors.New(message),
		StatusCode: http.StatusOK,
	}
}

func InternalServerError(message string) *HttpError {
	return &HttpError{
		Error:      errors.New(message),
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequestError(message string) *HttpError {
	return &HttpError{
		Error:      errors.New(message),
		StatusCode: http.StatusBadRequest,
	}
}
