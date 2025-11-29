package finance

import (
	"fmt"

	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

func (h *Handler) transfer(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransferRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Transfer(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully transfer\n================\nResult\nAccount: %v\nBalance: ฿%v", res.FromAccount, res.Balance), nil
}
