package finance

import (
	"fmt"
	"strings"

	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

func (h *Handler) getBalance() (string, *errors.AppError) {
	res, err := h.client.GetBalance()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Your balance\n\n")
	for _, v := range res.Accounts {
		sb.WriteString(fmt.Sprintf("Account: %v => Balance: ฿%v\n", v.Account, v.Balance))
	}
	return sb.String(), nil
}
