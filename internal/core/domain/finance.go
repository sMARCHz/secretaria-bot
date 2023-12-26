package domain

import (
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driven/financeservice/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Transaction
type TransactionRequest struct {
	Account     string  `json:"account"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
}

type TransactionResponse struct {
	Account string  `json:"account"`
	Balance float64 `json:"balance"`
}

func (t *TransactionRequest) ToProto() *pb.TransactionRequest {
	return &pb.TransactionRequest{
		AccountName: t.Account,
		Amount:      t.Amount,
		Category:    t.Category,
		Description: t.Description,
	}
}

// Transfer
type TransferRequest struct {
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
}

type TransferResponse struct {
	FromAccount string  `json:"from_account"`
	Balance     float64 `json:"balance"`
}

func (t *TransferRequest) ToProto() *pb.TransferRequest {
	return &pb.TransferRequest{
		FromAccountName: t.FromAccount,
		ToAccountName:   t.ToAccount,
		Amount:          t.Amount,
		Description:     t.Description,
	}
}

// GetBalance
type GetBalanceResponse struct {
	Accounts []AccountBalance `json:"accounts"`
}

type AccountBalance struct {
	Account string  `json:"account"`
	Balance float64 `json:"balance"`
}

// GetOverviewStatement
type GetOverviewStatementRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type GetOverviewStatementResponse struct {
	Revenue GetOverviewStatementSection `json:"revenue"`
	Expense GetOverviewStatementSection `json:"expense"`
	Profit  float64                     `json:"profit"`
}

type GetOverviewStatementSection struct {
	Total   float64            `json:"total"`
	Entries []CategorizedEntry `json:"entries"`
}

type CategorizedEntry struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

func (g *GetOverviewStatementRequest) ToProto() *pb.OverviewStatementRequest {
	return &pb.OverviewStatementRequest{
		From: timestamppb.New(g.From),
		To:   timestamppb.New(g.To),
	}
}
