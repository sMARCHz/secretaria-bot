package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	tokenizedMsg := []string{"!e", "debit1", "20000s"}
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().Deposit(&domain.TransactionRequest{
		Account:     "debit1",
		Amount:      20000,
		Category:    "s",
		Description: "",
	}).Return(&domain.TransactionResponse{
		Account: "debit1",
		Balance: 25000,
	}, nil)
	handler := NewHandler(client)

	res, err := handler.deposit(tokenizedMsg)

	assert.Nil(t, err)
	assert.Equal(t, "Succesfully deposit\n================\nResult\nAccount: debit1\nBalance: ฿25000", res)
	client.AssertExpectations(t)
}

func TestDeposit_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		mock         func(client *mocks.MockFinanceServiceClient)
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error when invalid command is provided",
			tokenizedMsg: []string{"!p"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!p/!e <account_name> <amount><category> <description>)"),
		},
		{
			it:           "return error when deposit fails",
			tokenizedMsg: []string{"!e", "debit1", "20000s"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Deposit(&domain.TransactionRequest{
					Account:     "debit1",
					Amount:      20000,
					Category:    "s",
					Description: "",
				}).Return(nil, errors.InternalServerError("failed to deposit"))
			},
			expectedErr: errors.InternalServerError("failed to deposit"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, err := handler.deposit(tc.tokenizedMsg)

			assert.Empty(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			client.AssertExpectations(t)
		})
	}
}
