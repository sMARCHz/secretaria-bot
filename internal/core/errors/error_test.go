package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := &AppError{Message: "An error occurred"}
	assert.Equal(t, "An error occurred", err.Error())
}

func TestBadRequestError(t *testing.T) {
	err := BadRequestError("Bad request")
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, "Bad request", err.Message)
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("Not found")
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "Not found", err.Message)
}

func TestUnprocessableEntityServerError(t *testing.T) {
	err := UnprocessableEntityServerError("Unprocessable entity")
	assert.Equal(t, http.StatusUnprocessableEntity, err.StatusCode)
	assert.Equal(t, "Unprocessable entity", err.Message)
}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("Internal server error")
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.Equal(t, "Internal server error", err.Message)
}

func TestBadGatewayError(t *testing.T) {
	err := BadGatewayError("Bad gateway")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
	assert.Equal(t, "Bad gateway", err.Message)
}
