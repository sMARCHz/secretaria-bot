package services

import (
	"net/http"
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/internal/core/services/finance"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewBotService(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)

	res := NewBotService(client)

	expected := &botServiceImpl{
		commandHandlers: []CommandHandler{
			finance.NewHandler(client),
		},
	}
	assert.Equal(t, expected, res)
}

func TestHandleTextMessage(t *testing.T) {
	testcases := []struct {
		it               string
		inputMsg         string
		expectedReplyMsg string
	}{
		{
			it:               "handle known finance command",
			inputMsg:         "balance",
			expectedReplyMsg: "Your balance\n\nAccount: debit1 => Balance: ฿1000\n",
		},
		{
			it:               "return 'No command input' for empty message",
			inputMsg:         "   ",
			expectedReplyMsg: "No command input",
		},
		{
			it:               "return 'Command not found' for unknown command",
			inputMsg:         "open sesame",
			expectedReplyMsg: "Command not found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			client := mocks.NewMockFinanceServiceClient(t)
			client.EXPECT().GetBalance().Return(&domain.GetBalanceResponse{
				Accounts: []domain.AccountBalance{
					{
						Account: "debit1",
						Balance: 1000,
					},
				},
			}, nil).Maybe()
			service := NewBotService(client)

			res, err := service.HandleTextMessage(tc.inputMsg)

			assert.Nil(t, err)
			assert.Equal(t, &domain.TextMessageResponse{
				ReplyMessage: tc.expectedReplyMsg,
			}, res)
			client.AssertExpectations(t)
		})
	}
}

func TestHandleTextMessage_Error(t *testing.T) {
	client := mocks.NewMockFinanceServiceClient(t)
	client.EXPECT().GetBalance().Return(nil, errors.InternalServerError("something went wrong")).Once()
	service := NewBotService(client)

	res, err := service.HandleTextMessage("balance")

	assert.Nil(t, res)
	assert.Equal(t, "something went wrong", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	client.AssertExpectations(t)
}
