package inbound

import (
	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
)

type BotService interface {
	HandleTextMessage(string) (*domain.TextMessageResponse, *errors.AppError)
}
