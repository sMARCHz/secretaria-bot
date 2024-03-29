package financetest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

type TestHandler struct {
	service services.BotService
}

func NewTestHandler(service services.BotService) TestHandler {
	return TestHandler{
		service: service,
	}
}

func (t *TestHandler) Test(ctx *gin.Context) {
	username, password, auth := ctx.Request.BasicAuth()
	if !auth || username != config.Get().App.TestUsername {
		logger.Warnf("someone tried to breach (username: %s, password: %s)", username, password)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var msg domain.TextMessageRequest
	if err := ctx.BindJSON(&msg); err != nil {
		logger.Error("cannot bind json: ", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := t.service.HandleTextMessage(msg.Message)
	if err != nil {
		ctx.AbortWithError(err.StatusCode, err)
	} else {
		ctx.JSON(http.StatusOK, res)
	}
}
