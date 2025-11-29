package line

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/secretaria-bot/internal/config"
	"github.com/sMARCHz/secretaria-bot/internal/logger"
	"github.com/sMARCHz/secretaria-bot/internal/ports/inbound"
)

type LineHandler struct {
	service inbound.BotService
	client  *linebot.Client
}

func NewLineHandler(service inbound.BotService) *LineHandler {
	lineCfg := config.Get().Line
	client, err := linebot.New(lineCfg.ChannelSecret, lineCfg.ChannelToken)
	if err != nil {
		logger.Fatal("cannot create linebot client: ", err)
	}
	return &LineHandler{
		service: service,
		client:  client,
	}
}

func (b *LineHandler) HandleLineMessage(ctx *gin.Context) {
	events, err := b.client.ParseRequest(ctx.Request)
	if err != nil {
		code := http.StatusInternalServerError
		if err == linebot.ErrInvalidSignature {
			code = http.StatusBadRequest
		}
		logger.Error("cannot parse line request: ", err)
		ctx.AbortWithError(code, err)
		return
	}

	b.processEvents(events)
}

func (b *LineHandler) processEvents(events []*linebot.Event) {
	for _, event := range events {
		if !isMyLineAccount(event) {
			b.replyMessage(event, "Unauthorized action!")
			continue
		}
		if event.Type != linebot.EventTypeMessage {
			continue
		}

		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			res, err := b.service.HandleTextMessage(message.Text)
			if err != nil {
				b.replyMessage(event, err.Message)
			} else {
				b.replyMessage(event, res.ReplyMessage)
			}
		default:
			b.replyMessage(event, "Unknown message type")
		}
	}
}

func (b *LineHandler) replyMessage(event *linebot.Event, replyMsg string) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
		logger.Error("cannot reply message: ", err)
	}
}

func isMyLineAccount(event *linebot.Event) bool {
	return event.Source.UserID == config.Get().Line.UserID
}
