package http

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/ports/inbound"
)

type testHandler struct {
	service inbound.BotService
}

func newTestHandler(service inbound.BotService) testHandler {
	return testHandler{
		service: service,
	}
}

type testMessage struct {
	Message string `json:"message"`
}

func (t *testHandler) handleTestMessage(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	var msg testMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	res, appErr := t.service.HandleTextMessage(msg.Message)
	if appErr != nil {
		ctx.AbortWithError(appErr.StatusCode, appErr)
		return
	}

	ctx.JSON(200, res.ReplyMessage)
}
