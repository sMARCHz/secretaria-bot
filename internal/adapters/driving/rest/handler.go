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
}

func (b *BotHandler) handleLineMessage(ctx *gin.Context) {
	bot, err := linebot.New(b.config.Line.ChannelSecret, b.config.Line.ChannelToken)
	if err != nil {
		b.logger.Error("Cannot create new linebot: ", err)
		return
	}

	events, err := bot.ParseRequest(ctx.Request)
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
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
					b.logger.Error("Cannot reply message: ", err)
				}
			}
		}
	}
}
