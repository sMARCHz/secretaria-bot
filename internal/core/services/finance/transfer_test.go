package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	tokenizedMsg := []string{"!t", "debit2", "debit1", "20000"}
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().Transfer(&domain.TransferRequest{
		FromAccount: "debit2",
		ToAccount:   "debit1",
		Amount:      20000,
		Description: "",
	}).Return(&domain.TransferResponse{
		FromAccount: "debit2",
		Balance:     500,
	}, nil)
	handler := NewHandler(client)

	res, err := handler.transfer(tokenizedMsg)

	assert.Nil(t, err)
	assert.Equal(t, "Succesfully transfer\n================\nResult\nAccount: debit2\nBalance: ฿500", res)
	client.AssertExpectations(t)
}

func TestTransfer_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		mock         func(client *mocks.MockFinanceServiceClient)
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error when invalid command is provided",
			tokenizedMsg: []string{"!t", "debit2"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!t <transfer_from> <transfer_to> <amount> <description>)"),
		},
		{
			it:           "return error when transfer fails",
			tokenizedMsg: []string{"!t", "debit2", "debit1", "20000"},
			mock: func(client *mocks.MockFinanceServiceClient) {
				client.EXPECT().Transfer(&domain.TransferRequest{
					FromAccount: "debit2",
					ToAccount:   "debit1",
					Amount:      20000,
					Description: "",
				}).Return(nil, errors.InternalServerError("failed to transfer"))
			},
			expectedErr: errors.InternalServerError("failed to transfer"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			if tc.mock != nil {
				tc.mock(client)
			}
			handler := NewHandler(client)

			res, err := handler.transfer(tc.tokenizedMsg)

			assert.Empty(t, res)
			assert.EqualError(t, err, tc.expectedErr.Message)
			client.AssertExpectations(t)
		})
	}
}
