package errs

import "net/http"

type AppError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
}

func NewNotFoundError(message string) *AppError {
	return &AppError{http.StatusNotFound, message}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{http.StatusInternalServerError, message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{http.StatusBadRequest, message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{http.StatusForbidden, message}
}
