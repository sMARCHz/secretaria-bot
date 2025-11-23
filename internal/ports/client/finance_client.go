package client

import (
	domain "github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

type FinanceServiceClient interface {
	Withdraw(*domain.TransactionRequest) (*domain.TransactionResponse, *errors.AppError)
	Deposit(*domain.TransactionRequest) (*domain.TransactionResponse, *errors.AppError)
	Transfer(*domain.TransferRequest) (*domain.TransferResponse, *errors.AppError)
	GetBalance() (*domain.GetBalanceResponse, *errors.AppError)
	GetOverviewStatement(*domain.GetOverviewStatementRequest) (*domain.GetOverviewStatementResponse, *errors.AppError)
	GetOverviewMonthlyStatement() (*domain.GetOverviewStatementResponse, *errors.AppError)
	GetOverviewAnnualStatement() (*domain.GetOverviewStatementResponse, *errors.AppError)
}
