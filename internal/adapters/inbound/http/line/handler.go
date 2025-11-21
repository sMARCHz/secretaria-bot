package line

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
	"github.com/sMARCHz/go-secretaria-bot/internal/ports/inbound"
)

type LineHandler struct {
	service inbound.BotService
}

func NewLineHandler(service inbound.BotService) LineHandler {
	return LineHandler{
		service: service,
	}
}

func (b *LineHandler) HandleLineMessage(ctx *gin.Context) {
	lineCfg := config.Get().Line
	client, err := linebot.New(lineCfg.ChannelSecret, lineCfg.ChannelToken)
	if err != nil {
		logger.Fatal("cannot create new linebot: ", err)
	}

	events, err := client.ParseRequest(ctx.Request)
	if err != nil {
		code := http.StatusInternalServerError
		if err == linebot.ErrInvalidSignature {
			code = http.StatusBadRequest
		}
		logger.Error("cannot parse line request: ", err)
		ctx.AbortWithError(code, err)
		return
	}

	b.processEvents(client, events)
}

func (b *LineHandler) processEvents(line *linebot.Client, events []*linebot.Event) {
	for _, event := range events {
		if !isMyLineAccount(event) {
			replyMessage(line, event, "Unauthorized action!")
			continue
		}
		if event.Type != linebot.EventTypeMessage {
			continue
		}

		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			res, err := b.service.HandleTextMessage(message.Text)
			if err != nil {
				replyMessage(line, event, err.Message)
			} else {
				replyMessage(line, event, res.ReplyMessage)
			}
		default:
			replyMessage(line, event, "Unknown message type")
		}
	}
}

func replyMessage(line *linebot.Client, event *linebot.Event, replyMsg string) {
	if _, err := line.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do(); err != nil {
		logger.Error("cannot reply message: ", err)
	}
}

func isMyLineAccount(event *linebot.Event) bool {
	return event.Source.UserID == config.Get().Line.UserID
}
