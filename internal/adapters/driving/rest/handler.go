package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

type Handler struct {
	service services.BotService
	config  config.Configuration
}

func (b *Handler) handleLineMessage(ctx *gin.Context) {
	line, err := linebot.New(b.config.Line.ChannelSecret, b.config.Line.ChannelToken)
	if err != nil {
		logger.Error("Cannot create new linebot: ", err)
	}

	events, err := line.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			logger.Error("Cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			logger.Error("Cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	for _, event := range events {
		if event.Source.UserID != b.config.Line.UserID {
			replyMessage(line, event, "Unauthorized action!")
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
				replyMessage(line, event, replyMsg)
			}
		}
	}
}

func replyMessage(line *linebot.Client, event *linebot.Event, replyMsg string) {
	if _, err := line.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
		logger.Error("Cannot reply message: ", err)
	}
}
