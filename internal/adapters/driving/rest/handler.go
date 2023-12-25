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
}

func NewHandler(service services.BotService) Handler {
	return Handler{
		service: service,
	}
}

func (b *Handler) handleLineMessage(ctx *gin.Context) {
	cfg := config.Get()
	line, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelToken)
	if err != nil {
		logger.Error("cannot create new linebot: ", err)
	}

	events, err := line.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			logger.Error("cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			logger.Error("cannot parse line request: ", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	for _, event := range events {
		if event.Source.UserID != cfg.Line.UserID {
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
		logger.Error("cannot reply message: ", err)
	}
}
