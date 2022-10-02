package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

type BotHandler struct {
	service services.BotService
	config  config.Configuration
	logger  logger.Logger
	linebot *linebot.Client
}

func (b *BotHandler) handleLineMessage(ctx *gin.Context) {
	events, err := b.linebot.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			b.logger.Error("Cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			b.logger.Error("Cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	for _, event := range events {
		if event.Source.UserID != b.config.Line.UserID {
			b.replyMessage(event, "Unauthorized action!")
			continue
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				replyMsg := ""
				res, appErr := b.service.HandleTextMessage(message.Text)
				if appErr != nil {
					replyMsg = appErr.Message
				} else {
					replyMsg = res.ReplyMessage
				}
				b.replyMessage(event, replyMsg)
			}
		}
	}
}

func (b *BotHandler) replyMessage(event *linebot.Event, replyMsg string) {
	if _, err := b.linebot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
		b.logger.Error("Cannot reply message: ", err)
	}
}
