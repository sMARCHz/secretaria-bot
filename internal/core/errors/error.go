package errors

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (a *AppError) Error() string {
	return a.Message
}

func BadRequestError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: msg}
}

func NotFoundError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: msg}
}

func UnprocessableEntityServerError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusUnprocessableEntity, Message: msg}
}

func InternalServerError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Message: msg}
}

func BadGatewayError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadGateway, Message: msg}
}
