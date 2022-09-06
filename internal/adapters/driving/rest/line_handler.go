package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func handleLineMessage(ctx *gin.Context, config config.LineMessageConfiguration, logger logger.Logger) {
	bot, err := linebot.New(config.ChannelSecret, config.ChannelToken)
	if err != nil {
		logger.Error("Cannot create new linebot: ", err)
		return
	}

	events, err := bot.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.AbortWithError(http.StatusBadRequest, err)
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					logger.Error(err)
				}
			}
		}
	}
}
