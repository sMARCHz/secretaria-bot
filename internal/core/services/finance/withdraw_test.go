package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWithdraw(t *testing.T) {
	tokenizedMsg := []string{"!p", "debit1", "500sh", "youtube membership"}
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().Withdraw(&domain.TransactionRequest{
		Account:     "debit1",
		Amount:      500,
		Category:    "sh",
		Description: "youtube membership",
	}).Return(&domain.TransactionResponse{
		Account: "debit1",
		Balance: 1000,
	}, nil)
	handler := NewHandler(client)

	res, err := handler.withdraw(tokenizedMsg)

	assert.Nil(t, err)
	assert.Equal(t, "Succesfully withdraw\n================\nResult\nAccount: debit1\nBalance: ฿1000", res)
	client.AssertExpectations(t)
}

func TestWithdraw_Error(t *testing.T) {
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
			it:           "return error when withdraw fails",
			tokenizedMsg: []string{"!p", "debit1", "500sh", "youtube membership"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Withdraw(&domain.TransactionRequest{
					Account:     "debit1",
					Amount:      500,
					Category:    "sh",
					Description: "youtube membership",
				}).Return(nil, errors.InternalServerError("failed to withdraw"))
			},
			expectedErr: errors.InternalServerError("failed to withdraw"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, err := handler.withdraw(tc.tokenizedMsg)

			assert.Empty(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			client.AssertExpectations(t)
		})
	}
}
