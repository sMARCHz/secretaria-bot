package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().GetBalance().Return(&domain.GetBalanceResponse{
		Accounts: []domain.AccountBalance{
			{
				Account: "debit1",
				Balance: 5000,
			},
		},
	}, nil)
	handler := NewHandler(client)

	res, err := handler.getBalance()

	assert.Nil(t, err)
	assert.Equal(t, "Your balance\n\nAccount: debit1 => Balance: ฿5000\n", res)
	client.AssertExpectations(t)
}

func TestGetBalance_Error(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().GetBalance().Return(nil, errors.InternalServerError("something went wrong"))
	handler := NewHandler(client)

	res, err := handler.getBalance()

	assert.Empty(t, res)
	assert.EqualError(t, err, "something went wrong")
	client.AssertExpectations(t)
}
