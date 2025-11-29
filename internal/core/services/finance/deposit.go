package finance

import (
	"fmt"

	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

func (h *Handler) deposit(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransactionRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Deposit(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully deposit\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}
