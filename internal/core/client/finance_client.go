package client

import (
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/core/dto"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
)

type FinanceServiceClient interface {
	Withdraw([]string) (*dto.TransactionResponse, *errors.AppError)
	Deposit([]string) (*dto.TransactionResponse, *errors.AppError)
	Transfer([]string) (*dto.TransferResponse, *errors.AppError)
	GetBalance() (*dto.GetBalanceResponse, *errors.AppError)
	GetOverviewStatement(time.Time, time.Time) (*dto.GetOverviewStatementResponse, *errors.AppError)
	GetOverviewMonthlyStatement() (*dto.GetOverviewStatementResponse, *errors.AppError)
	GetOverviewAnnualStatement() (*dto.GetOverviewStatementResponse, *errors.AppError)
}
