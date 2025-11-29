package http

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	apperrors "github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewTestHandler(t *testing.T) {
	bot := mocks.NewMockBotService(t)
	handler := newTestHandler(bot)
	assert.Equal(t, &testHandler{service: bot}, handler)
}

func TestHandleTextMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := mocks.NewMockBotService(t)
	bot.EXPECT().HandleTextMessage("hello").Return(&domain.TextMessageResponse{
		ReplyMessage: "world",
	}, nil)

	requestBody := `{"message":"hello"}`
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/__test", strings.NewReader(requestBody))

	handler := &testHandler{service: bot}

	handler.handleTestMessage(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"world"}`, w.Body.String())
	bot.AssertExpectations(t)
}

func TestHandleTextMessage_Error(t *testing.T) {
	testcases := []struct {
		it                 string
		body               io.Reader
		mock               func(bot *mocks.MockBotService)
		expectedHTTPStatus int
		expectedBody       string
	}{
		{
			it:                 "returns error with status 500 when fails to read request body",
			body:               errorReader{},
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedBody:       `{"error":"read error"}`,
		},
		{
			it:                 "returns error with status 500 when fails to parse request body",
			body:               strings.NewReader(`{"message":}`),
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedBody:       `{"error":"invalid character '}' looking for beginning of value"}`,
		},
		{
			it:   "returns error with status and message from service layer when fails to handle the message",
			body: strings.NewReader(`{"message":"hello"}`),
			mock: func(bot *mocks.MockBotService) {
				bot.EXPECT().HandleTextMessage("hello").Return(nil, apperrors.BadGatewayError("fail to handle message"))
			},
			expectedHTTPStatus: http.StatusBadGateway,
			expectedBody:       `{"error":"fail to handle message"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			bot := mocks.NewMockBotService(t)
			if tc.mock != nil {
				tc.mock(bot)
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("POST", "/__test", tc.body)

			handler := &testHandler{service: bot}

			handler.handleTestMessage(ctx)

			assert.Equal(t, tc.expectedHTTPStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
			bot.AssertExpectations(t)
		})
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) { return 0, errors.New("read error") }
