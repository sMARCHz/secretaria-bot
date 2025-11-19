package services

import (
	"strings"

	"github.com/sMARCHz/go-secretaria-bot/internal/core/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
)

type BotService interface {
	HandleTextMessage(string) (*domain.TextMessageResponse, *errors.AppError)
}

type botServiceImpl struct {
	financeClient client.FinanceServiceClient
}

func NewBotService(financeClient client.FinanceServiceClient) BotService {
	return &botServiceImpl{
		financeClient: financeClient,
	}
}

func (b *botServiceImpl) HandleTextMessage(msg string) (*domain.TextMessageResponse, *errors.AppError) {
	replyMsg := ""
	msg = strings.TrimSpace(msg)
	msg = strings.ToLower(msg)
	msgArr := strings.Fields(msg)

	var err *errors.AppError
	switch msgArr[0] {
	case "!p":
		replyMsg, err = b.callWithdraw(msgArr)
	case "!e":
		replyMsg, err = b.callDeposit(msgArr)
	case "!t":
		replyMsg, err = b.callTransfer(msgArr)
	case "balance":
		replyMsg, err = b.callBalance()
	case "statement":
		replyMsg, err = b.callStatement(msgArr)
	default:
		replyMsg = "Command not found"
	}

	if err != nil {
		return nil, err
	}
	return &domain.TextMessageResponse{ReplyMessage: replyMsg}, nil
}
