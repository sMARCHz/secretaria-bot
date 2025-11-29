package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/secretaria-bot/internal/ports/inbound"
)

type testHandler struct {
	service inbound.BotService
}

func newTestHandler(service inbound.BotService) *testHandler {
	return &testHandler{
		service: service,
	}
}

type testMessage struct {
	Message string `json:"message"`
}

func (t *testHandler) handleTestMessage(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var msg testMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, appErr := t.service.HandleTextMessage(msg.Message)
	if appErr != nil {
		ctx.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	ctx.JSON(200, gin.H{"message": res.ReplyMessage})
}
