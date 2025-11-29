package finance

import (
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/internal/ports/client"
)

// Handler implements command handling for finance-related commands.
type Handler struct {
	client client.FinanceServiceClient
}

// NewHandler constructs a finance command handler.
func NewHandler(client client.FinanceServiceClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Match(cmd string) bool {
	if _, exist := commandPrefixSet[cmd]; exist {
		return true
	}
	return false
}

func (h *Handler) Handle(tokenizedMsg []string) (string, *errors.AppError) {
	if len(tokenizedMsg) == 0 {
		return "", errors.BadRequestError(commandNotFoundMsg)
	}
	switch tokenizedMsg[0] {
	case "!p":
		return h.withdraw(tokenizedMsg)
	case "!e":
		return h.deposit(tokenizedMsg)
	case "!t":
		return h.transfer(tokenizedMsg)
	case "balance":
		return h.getBalance()
	case "statement":
		return h.getStatement(tokenizedMsg)
	default:
		return "", errors.BadRequestError(invalidCommandMsg)
	}
}
