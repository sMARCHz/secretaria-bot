package finance

import (
	"fmt"

	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

func (h *Handler) withdraw(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransactionRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Withdraw(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully withdraw\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}
