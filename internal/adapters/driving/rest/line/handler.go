package line

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

type LineHandler struct {
	service services.BotService
}

func NewLineHandler(service services.BotService) LineHandler {
	return LineHandler{
		service: service,
	}
}

func (b *LineHandler) HandleLineMessage(ctx *gin.Context) {
	cfg := config.Get()
	line, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelToken)
	if err != nil {
		logger.Fatal("cannot create new linebot: ", err)
	}

	events, err := line.ParseRequest(ctx.Request)
	if err != nil {
		code := http.StatusInternalServerError
		if err == linebot.ErrInvalidSignature {
			code = http.StatusBadRequest
		}
		logger.Error("cannot parse line request: ", err)
		ctx.AbortWithError(code, err)
		return
	}

	b.processEvents(line, events)
}

func (b *LineHandler) processEvents(line *linebot.Client, events []*linebot.Event) {
	for _, event := range events {
		if !isMe(event) {
			replyMessage(line, event, "Unauthorized action!")
			continue
		}
		if event.Type != linebot.EventTypeMessage {
			continue
		}

		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			res, err := b.service.HandleTextMessage(message.Text)
			replyMsg := res.ReplyMessage
			if err != nil {
				replyMsg = errors.GetErrorMessage(err)
			}
			replyMessage(line, event, replyMsg)
		}
	}
}

func replyMessage(line *linebot.Client, event *linebot.Event, replyMsg string) {
	if _, err := line.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
		logger.Error("cannot reply message: ", err)
	}
}

func isMe(event *linebot.Event) bool {
	return event.Source.UserID == config.Get().Line.UserID
}
