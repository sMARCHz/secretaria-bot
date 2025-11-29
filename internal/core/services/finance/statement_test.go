package finance

import (
	"net/http"
	"testing"
	"time"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetStatement(t *testing.T) {
	testcases := []struct {
		it               string
		tokenizedMsg     []string
		mock             func(client *mocks.MockFinanceServiceClient)
		expectedReplyMsg string
	}{
		{
			it:           "return reply message for monthly statement if no argument is provided",
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
		{
			it:           "return reply message for annual statement if 'a' argument is provided",
			tokenizedMsg: []string{"statement", "a"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewAnnualStatement().Return(&domain.GetOverviewStatementResponse{
					Revenue: &domain.GetOverviewStatementSection{
						Total: 240000,
					},
					Expense: &domain.GetOverviewStatementSection{
						Total: 180000,
					},
					Profit: 60000,
				}, nil)
			},
			expectedReplyMsg: "Annual Statement\n================\nRevenue: ฿240000\n\nExpense: ฿180000\n\nProfit: ฿60000",
		},
		{
			it:           "return reply message for selected range statement if two date arguments are provided",
			tokenizedMsg: []string{"statement", "2025-01-01", "2025-03-31"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewStatement(&domain.GetOverviewStatementRequest{
					From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
				}).Return(&domain.GetOverviewStatementResponse{
					Revenue: &domain.GetOverviewStatementSection{
						Total: 60000,
					},
					Expense: &domain.GetOverviewStatementSection{
						Total: 45000,
					},
					Profit: 15000,
				}, nil)
			},
			expectedReplyMsg: "Income Statement\n================\nRevenue: ฿60000\n\nExpense: ฿45000\n\nProfit: ฿15000",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			tc.mock(client)
			handler := NewHandler(client)

			res, err := handler.getStatement(tc.tokenizedMsg)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedReplyMsg, res)
			client.AssertExpectations(t)
		})
	}
}

func TestGetStatement_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		mock         func(client *mocks.MockFinanceServiceClient)
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error when invalid command is provided",
			tokenizedMsg: []string{"this", "is", "invalid", "command"},
			expectedErr:  errors.BadRequestError("Invalid command"),
		},
		{
			it:           "return error when fail to get statement",
			tokenizedMsg: []string{"statement"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewMonthlyStatement().Return(nil, errors.InternalServerError("failed to get statement"))
			},
			expectedErr: errors.InternalServerError("failed to get statement"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, err := handler.getStatement(tc.tokenizedMsg)

			assert.Empty(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			assert.Equal(t, tc.expectedErr.StatusCode, err.StatusCode)
		})
	}
}

func TestCallMonthlyOrAnnualStatement(t *testing.T) {
	testcases := []struct {
		it                    string
		statementType         string
		mock                  func(client *mocks.MockFinanceServiceClient)
		expectedRes           *domain.GetOverviewStatementResponse
		expectedStatementType string
	}{
		{
			it:            "return monthly statement when statementType is m",
			statementType: "m",
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewMonthlyStatement().Return(&domain.GetOverviewStatementResponse{
					Revenue: &domain.GetOverviewStatementSection{
						Total: 20000,
						Entries: []domain.CategorizedEntry{
							{Category: "Salary", Amount: 20000},
						},
					},
					Expense: &domain.GetOverviewStatementSection{
						Total: 15000,
						Entries: []domain.CategorizedEntry{
							{Category: "Food", Amount: 5000},
						},
					},
				}, nil)
			},
			expectedRes: &domain.GetOverviewStatementResponse{
				Revenue: &domain.GetOverviewStatementSection{
					Total: 20000,
					Entries: []domain.CategorizedEntry{
						{Category: "Salary", Amount: 20000},
					},
				},
				Expense: &domain.GetOverviewStatementSection{
					Total: 15000,
					Entries: []domain.CategorizedEntry{
						{Category: "Food", Amount: 5000},
					},
				},
			},
			expectedStatementType: "Monthly",
		},
		{
			it:            "return annual statement when statementType is a",
			statementType: "a",
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewAnnualStatement().Return(&domain.GetOverviewStatementResponse{
					Revenue: &domain.GetOverviewStatementSection{
						Total: 240000,
						Entries: []domain.CategorizedEntry{
							{Category: "Salary", Amount: 240000},
						},
					},
					Expense: &domain.GetOverviewStatementSection{
						Total: 180000,
						Entries: []domain.CategorizedEntry{
							{Category: "Food", Amount: 60000},
							{Category: "Shopping", Amount: 120000},
						},
					},
				}, nil)
			},
			expectedRes: &domain.GetOverviewStatementResponse{
				Revenue: &domain.GetOverviewStatementSection{
					Total: 240000,
					Entries: []domain.CategorizedEntry{
						{Category: "Salary", Amount: 240000},
					},
				},
				Expense: &domain.GetOverviewStatementSection{
					Total: 180000,
					Entries: []domain.CategorizedEntry{
						{Category: "Food", Amount: 60000},
						{Category: "Shopping", Amount: 120000},
					},
				},
			},
			expectedStatementType: "Annual",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, statementType, err := handler.callMonthlyOrAnnualStatement(tc.statementType)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedRes, res)
			assert.Equal(t, tc.expectedStatementType, statementType)
		})
	}
}

func TestCallMonthlyOrAnnualStatement_Error(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)
	handler := NewHandler(client)

	res, statementType, err := handler.callMonthlyOrAnnualStatement("invalid_type")

	assert.Nil(t, res)
	assert.Empty(t, statementType)
	assert.EqualError(t, err, "Invalid command")
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestCallSelectedRangeStatement(t *testing.T) {
	financeRes := &domain.GetOverviewStatementResponse{
		Revenue: &domain.GetOverviewStatementSection{
			Total: 220000,
		},
		Expense: &domain.GetOverviewStatementSection{
			Total: 160000,
		},
		Profit: 60000,
	}
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().GetOverviewStatement(&domain.GetOverviewStatementRequest{
		From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2025, 11, 23, 0, 0, 0, 0, time.UTC),
	}).Return(financeRes, nil)
	handler := NewHandler(client)

	res, err := handler.callSelectedRangeStatement("2025-01-01", "2025-11-23")

	assert.Nil(t, err)
	assert.Equal(t, financeRes, res)
}

func TestCallSelectedRangeStatement_Error(t *testing.T) {
	testcases := []struct {
		it          string
		from        string
		to          string
		mock        func(client *mocks.MockFinanceServiceClient)
		expectedErr *errors.AppError
	}{
		{
			it:          "return error when fail to parse from date",
			from:        "invalid-date",
			to:          "2025-12-31",
			expectedErr: errors.BadRequestError("Invalid command's arguments.\nPlease recheck the from_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>"),
		},
		{
			it:          "return error when fail to parse to date",
			from:        "2025-01-01",
			to:          "invalid-date",
			expectedErr: errors.BadRequestError("Invalid command's arguments.\nPlease recheck the to_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>"),
		},
		{
			it:   "return error when fail to get statement from finance service",
			from: "2025-01-01",
			to:   "2025-12-31",
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().GetOverviewStatement(&domain.GetOverviewStatementRequest{
					From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				}).Return(nil, errors.InternalServerError("failed to get statement"))
			},
			expectedErr: errors.InternalServerError("failed to get statement"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, err := handler.callSelectedRangeStatement(tc.from, tc.to)

			assert.Nil(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			assert.Equal(t, tc.expectedErr.StatusCode, err.StatusCode)
		})
	}
}

func TestPrintStatement(t *testing.T) {
	testcases := []struct {
		it           string
		statementRes *domain.GetOverviewStatementResponse
		expected     string
	}{
		{
			it: "returns string statement",
			statementRes: &domain.GetOverviewStatementResponse{
				Revenue: &domain.GetOverviewStatementSection{
					Total: 30000,
					Entries: []domain.CategorizedEntry{
						{Category: "Salary", Amount: 30000},
					},
				},
				Expense: &domain.GetOverviewStatementSection{
					Total: 20000,
					Entries: []domain.CategorizedEntry{
						{Category: "Food", Amount: 8000},
						{Category: "Shopping", Amount: 12000},
					},
				},
				Profit: 10000,
			},
			expected: "Income Statement\n================\nRevenue: ฿30000\nSalary = ฿30000\n\nExpense: ฿20000\nFood = ฿8000\nShopping = ฿12000\n\nProfit: ฿10000",
		},
		{
			it: "returns string statement with revenue=0 and no entries when the revenue object isn't in the response",
			statementRes: &domain.GetOverviewStatementResponse{
				Expense: &domain.GetOverviewStatementSection{
					Total: 20000,
					Entries: []domain.CategorizedEntry{
						{Category: "Food", Amount: 8000},
						{Category: "Shopping", Amount: 12000},
					},
				},
				Profit: -20000,
			},
			expected: "Income Statement\n================\nRevenue: ฿0\n\nExpense: ฿20000\nFood = ฿8000\nShopping = ฿12000\n\nProfit: ฿-20000",
		},
		{
			it: "returns string statement with expense=0 and no entries when the expense object isn't in the response",
			statementRes: &domain.GetOverviewStatementResponse{
				Revenue: &domain.GetOverviewStatementSection{
					Total: 30000,
					Entries: []domain.CategorizedEntry{
						{Category: "Salary", Amount: 30000},
					},
				},
				Profit: 30000,
			},
			expected: "Income Statement\n================\nRevenue: ฿30000\nSalary = ฿30000\n\nExpense: ฿0\n\nProfit: ฿30000",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			msg := printStatement(tc.statementRes, "Income")
			assert.Equal(t, tc.expected, msg)
		})
	}
}
