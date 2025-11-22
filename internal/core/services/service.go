package services

import (
	"strings"

	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services/finance"
	"github.com/sMARCHz/go-secretaria-bot/internal/ports/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/ports/inbound"
)

type botServiceImpl struct {
	commandHandlers []CommandHandler
}

func NewBotService(financeClient client.FinanceServiceClient) inbound.BotService {
	financeHandler := finance.NewHandler(financeClient)
	return &botServiceImpl{
		commandHandlers: []CommandHandler{financeHandler},
	}
}

func (b *botServiceImpl) HandleTextMessage(msg string) (*domain.TextMessageResponse, *errors.AppError) {
	msg = strings.TrimSpace(msg)
	msg = strings.ToLower(msg)
	tokenizedMsg := strings.Fields(msg)
	if len(tokenizedMsg) == 0 {
		return &domain.TextMessageResponse{ReplyMessage: "No command input"}, nil
	}

	var err *errors.AppError
	var handled bool
	var replyMsg string
	for _, h := range b.commandHandlers {
		if h.Match(tokenizedMsg[0]) {
			replyMsg, err = h.Handle(tokenizedMsg)
			handled = true
			break
		}
	}
	if err != nil {
		return nil, err
	}
	if !handled {
		replyMsg = "Command not found"
	}

	return &domain.TextMessageResponse{
		ReplyMessage: replyMsg,
	}, nil
}
