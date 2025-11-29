package domain

import (
	"testing"
	"time"

	"github.com/sMARCHz/secretaria-bot/internal/adapters/client/finance/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransactionRequestToProto(t *testing.T) {
	req := &TransactionRequest{
		Account:     "debit1",
		Amount:      100.50,
		Category:    "sh",
		Description: "Weekly shopping",
	}

	res := req.ToProto()

	expected := &pb.TransactionRequest{
		AccountName: "debit1",
		Amount:      100.50,
		Category:    "sh",
		Description: "Weekly shopping",
	}
	assert.Equal(t, expected, res)
}

func TestTransferRequestToProto(t *testing.T) {
	req := &TransferRequest{
		FromAccount: "debit2",
		ToAccount:   "debit1",
		Amount:      250.00,
		Description: "Monthly transfer",
	}

	res := req.ToProto()

	expected := &pb.TransferRequest{
		FromAccountName: "debit2",
		ToAccountName:   "debit1",
		Amount:          250.00,
		Description:     "Monthly transfer",
	}
	assert.Equal(t, expected, res)
}

func TestGetOverviewStatementRequestToProto(t *testing.T) {
	req := &GetOverviewStatementRequest{
		From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
	}

	res := req.ToProto()

	expected := &pb.OverviewStatementRequest{
		From: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		To:   timestamppb.New(time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)),
	}
	assert.Equal(t, expected, res)

}
