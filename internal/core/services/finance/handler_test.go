package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)

	res := NewHandler(client)

	expected := &Handler{client: client}
	assert.Equal(t, expected, res)
}

func TestMatch(t *testing.T) {
	testcases := []struct {
		it       string
		cmd      string
		expected bool
	}{
		{
			it:       "return true for valid command",
			cmd:      "!p",
			expected: true,
		},
		{
			it:       "return false for invalid or unknown command",
			cmd:      "unknown",
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			handler := NewHandler(client)

			res := handler.Match(tc.cmd)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestHandle(t *testing.T) {
	testcases := []struct {
		it               string
		tokenizedMsg     []string
		mock             func(client *mocks.MockFinanceServiceClient)
		expectedReplyMsg string
	}{
		{
			it:           "return reply message for withdraw command",
			tokenizedMsg: []string{"!p", "debit1", "500sh", "youtube membership"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Withdraw(&domain.TransactionRequest{
					Account:     "debit1",
					Amount:      500,
					Category:    "sh",
					Description: "youtube membership",
				}).Return(&domain.TransactionResponse{
					Account: "debit1",
					Balance: 1000,
				}, nil)
			},
			expectedReplyMsg: "Succesfully withdraw\n================\nResult\nAccount: debit1\nBalance: ฿1000",
		},
		{
			it:           "return reply message for deposit command",
			tokenizedMsg: []string{"!e", "debit1", "20000s"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Deposit(&domain.TransactionRequest{
					Account:     "debit1",
					Amount:      20000,
					Category:    "s",
					Description: "",
				}).Return(&domain.TransactionResponse{
					Account: "debit1",
					Balance: 25000,
				}, nil)
			},
			expectedReplyMsg: "Succesfully deposit\n================\nResult\nAccount: debit1\nBalance: ฿25000",
		},
		{
			it:           "return reply message for transfer command",
			tokenizedMsg: []string{"!t", "debit2", "debit1", "20000"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Transfer(&domain.TransferRequest{
					FromAccount: "debit2",
					ToAccount:   "debit1",
					Amount:      20000,
					Description: "",
				}).Return(&domain.TransferResponse{
					FromAccount: "debit2",
					Balance:     500,
				}, nil)
			},
			expectedReplyMsg: "Succesfully transfer\n================\nResult\nAccount: debit2\nBalance: ฿500",
		},
		{
			it:           "return reply message for balance command",
			tokenizedMsg: []string{"balance"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetBalance().Return(&domain.GetBalanceResponse{
					Accounts: []domain.AccountBalance{
						{
							Account: "debit1",
							Balance: 5000,
						},
					},
				}, nil)
			},
			expectedReplyMsg: "Your balance\n\nAccount: debit1 => Balance: ฿5000\n",
		},
		{
			it:           "return reply message for statement command",
			tokenizedMsg: []string{"statement"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewMonthlyStatement().Return(&domain.GetOverviewStatementResponse{
					Revenue: &domain.GetOverviewStatementSection{
						Total: 20000,
					},
					Expense: &domain.GetOverviewStatementSection{
						Total: 15000,
					},
					Profit: 5000,
				}, nil)
			},
			expectedReplyMsg: "Monthly Statement\n================\nRevenue: ฿20000\n\nExpense: ฿15000\n\nProfit: ฿5000",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			tc.mock(client)
			handler := NewHandler(client)

			replyMsg, err := handler.Handle(tc.tokenizedMsg)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedReplyMsg, replyMsg)
			client.AssertExpectations(t)
		})
	}
}

func TestHandle_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error bad request if tokenizedMsg is an empty slice",
			tokenizedMsg: []string{},
			expectedErr:  errors.BadRequestError("Command not found"),
		},
		{
			it:           "return error from handler method",
			tokenizedMsg: []string{"!p", "invalid"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!p/!e <account_name> <amount><category> <description>)"),
		},
		{
			it:           "return error for unknown command",
			tokenizedMsg: []string{"unknown"},
			expectedErr:  errors.BadRequestError("Invalid command"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			handler := NewHandler(client)

			res, err := handler.Handle(tc.tokenizedMsg)

			assert.Empty(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			assert.Equal(t, tc.expectedErr.StatusCode, err.StatusCode)
		})
	}
}
