package finance

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/adapters/client/finance/pb"
	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewFinanceServiceClient(t *testing.T) {
	t.Skip("integration test - requires running finance gRPC service; skipped in unit tests")
}

func TestWithdraw(t *testing.T) {
	gRPCRes := &pb.TransactionResponse{
		AccountName: "debit1",
		Balance:     500,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Withdraw", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Withdraw(&domain.TransactionRequest{})

	expected := &domain.TransactionResponse{
		Account: gRPCRes.AccountName,
		Balance: gRPCRes.Balance,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestWithdraw_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Withdraw", mock.Anything, mock.Anything).Return(nil, errors.New("fails to withdraw"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Withdraw(&domain.TransactionRequest{})

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot withdraw money: fails to withdraw")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestDeposit(t *testing.T) {
	gRPCRes := &pb.TransactionResponse{
		AccountName: "debit1",
		Balance:     500,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Deposit", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Deposit(&domain.TransactionRequest{})

	expected := &domain.TransactionResponse{
		Account: gRPCRes.AccountName,
		Balance: gRPCRes.Balance,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestDeposit_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Deposit", mock.Anything, mock.Anything).Return(nil, errors.New("fails to deposit"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Deposit(&domain.TransactionRequest{})

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot deposit money: fails to deposit")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestTransfer(t *testing.T) {
	gRPCRes := &pb.TransferResponse{
		FromAccountName: "debit1",
		Balance:         500,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Transfer", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Transfer(&domain.TransferRequest{})

	expected := &domain.TransferResponse{
		FromAccount: gRPCRes.FromAccountName,
		Balance:     gRPCRes.Balance,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestTransfer_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("Transfer", mock.Anything, mock.Anything).Return(nil, errors.New("fails to transfer"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.Transfer(&domain.TransferRequest{})

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot transfer money: fails to transfer")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestGetBalance(t *testing.T) {
	gRPCRes := &pb.GetBalanceResponse{
		Accounts: []*pb.AccountBalance{
			{
				AccountName: "debit1",
				Balance:     500,
			},
		},
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetBalance", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetBalance()

	expected := &domain.GetBalanceResponse{
		Accounts: []domain.AccountBalance{
			{
				Account: "debit1",
				Balance: 500,
			},
		},
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestGetBalance_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetBalance", mock.Anything, mock.Anything).Return(nil, errors.New("fails to get balance"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetBalance()

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot get balance: fails to get balance")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestGetOverviewStatement(t *testing.T) {
	gRPCRes := &pb.OverviewStatementResponse{
		Revenue: &pb.OverviewStatementSection{
			Total: 20000,
		},
		Expense: &pb.OverviewStatementSection{
			Total: 10000,
		},
		Profit: 10000,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewStatement", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewStatement(&domain.GetOverviewStatementRequest{})

	expected := &domain.GetOverviewStatementResponse{
		Revenue: res.Revenue,
		Expense: res.Expense,
		Profit:  res.Profit,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOverviewStatement_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewStatement", mock.Anything, mock.Anything).Return(nil, errors.New("fails to get overview statement"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewStatement(&domain.GetOverviewStatementRequest{})

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot get overview statement: fails to get overview statement")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestGetOverviewMonthlyStatement(t *testing.T) {
	gRPCRes := &pb.OverviewStatementResponse{
		Revenue: &pb.OverviewStatementSection{
			Total: 20000,
		},
		Expense: &pb.OverviewStatementSection{
			Total: 10000,
		},
		Profit: 10000,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewMonthlyStatement", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewMonthlyStatement()

	expected := &domain.GetOverviewStatementResponse{
		Revenue: res.Revenue,
		Expense: res.Expense,
		Profit:  res.Profit,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOverviewMonthlyStatement_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewMonthlyStatement", mock.Anything, mock.Anything).Return(nil, errors.New("fails to get monthly statement"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewMonthlyStatement()

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot get monthly overview statement: fails to get monthly statement")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestGetOverviewAnnualStatement(t *testing.T) {
	gRPCRes := &pb.OverviewStatementResponse{
		Revenue: &pb.OverviewStatementSection{
			Total: 20000,
		},
		Expense: &pb.OverviewStatementSection{
			Total: 10000,
		},
		Profit: 10000,
	}
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewAnnualStatement", mock.Anything, mock.Anything).Return(gRPCRes, nil)
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewAnnualStatement()

	expected := &domain.GetOverviewStatementResponse{
		Revenue: res.Revenue,
		Expense: res.Expense,
		Profit:  res.Profit,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOverviewAnnualStatement_Error(t *testing.T) {
	gRPCClient := mocks.NewMockGRPCFinanceServiceClient(t)
	gRPCClient.On("GetOverviewAnnualStatement", mock.Anything, mock.Anything).Return(nil, errors.New("fails to get annual statement"))
	client := &financeServiceClient{
		client: gRPCClient,
	}

	res, err := client.GetOverviewAnnualStatement()

	assert.Nil(t, res)
	assert.EqualError(t, err, "cannot get annual overview statement: fails to get annual statement")
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestToGetOverviewStatementResponse(t *testing.T) {
	testcases := []struct {
		it       string
		gRPCRes  *pb.OverviewStatementResponse
		expected *domain.GetOverviewStatementResponse
	}{
		{
			it:       "returns empty response if gRPC response is nil",
			gRPCRes:  nil,
			expected: &domain.GetOverviewStatementResponse{},
		},
		{
			it: "returns response with revenue if gRPC response contains revenue object",
			gRPCRes: &pb.OverviewStatementResponse{
				Revenue: &pb.OverviewStatementSection{
					Total: 15000,
					Entries: []*pb.CategorizedEntry{
						{
							Category: "s",
							Amount:   10000,
						},
						{
							Category: "misc",
							Amount:   5000,
						},
					},
				},
				Profit: 15000,
			},
			expected: &domain.GetOverviewStatementResponse{
				Revenue: &domain.GetOverviewStatementSection{
					Total: 15000,
					Entries: []domain.CategorizedEntry{
						{
							Category: "s",
							Amount:   10000,
						},
						{
							Category: "misc",
							Amount:   5000,
						},
					},
				},
				Expense: nil,
				Profit:  15000,
			},
		},
		{
			it: "returns response with expense if gRPC response contains expense object",
			gRPCRes: &pb.OverviewStatementResponse{
				Expense: &pb.OverviewStatementSection{
					Total: 1000,
					Entries: []*pb.CategorizedEntry{
						{
							Category: "sh",
							Amount:   800,
						},
						{
							Category: "sn",
							Amount:   200,
						},
					},
				},
				Profit: -1000,
			},
			expected: &domain.GetOverviewStatementResponse{
				Revenue: nil,
				Expense: &domain.GetOverviewStatementSection{
					Total: 1000,
					Entries: []domain.CategorizedEntry{
						{
							Category: "sh",
							Amount:   800,
						},
						{
							Category: "sn",
							Amount:   200,
						},
					},
				},
				Profit: -1000,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := &financeServiceClient{}
			res := client.toGetOverviewStatementResponse(tc.gRPCRes)

			assert.Equal(t, tc.expected, res)
		})
	}
}
