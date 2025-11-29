package finance

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseTransactionRequest(t *testing.T) {
	tokenizedMsg := []string{"!p", "debit1", "200sh", "steam purchase"}

	res, err := parseTransactionRequest(tokenizedMsg)

	expected := &domain.TransactionRequest{
		Account:     "debit1",
		Amount:      200,
		Category:    "sh",
		Description: "steam purchase",
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestParseTransactionRequest_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error when command's length is less than 3",
			tokenizedMsg: []string{"!p", "debit1"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!p/!e <account_name> <amount><category> <description>)"),
		},
		{
			it:           "return error when amount and category combination is invalid",
			tokenizedMsg: []string{"!p", "debit1", "200"},
			expectedErr:  errors.BadRequestError("Invalid amount and category combination"),
		},
		{
			it:           "return error when amount cannot be parsed to float64",
			tokenizedMsg: []string{"!p", "debit1", "sh"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			res, err := parseTransactionRequest(tc.tokenizedMsg)
			assert.Nil(t, res)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedErr.StatusCode, err.StatusCode)
		})
	}
}

func TestParseTransferRequest(t *testing.T) {
	tokenizedMsg := []string{"!t", "debit2", "debit1", "20000", "salary"}

	res, err := parseTransferRequest(tokenizedMsg)

	expected := &domain.TransferRequest{
		FromAccount: "debit2",
		ToAccount:   "debit1",
		Amount:      20000,
		Description: "salary",
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestParseTransferRequest_Error(t *testing.T) {
	testcases := []struct {
		it           string
		tokenizedMsg []string
		expectedErr  *errors.AppError
	}{
		{
			it:           "return error when command's length is less than 4",
			tokenizedMsg: []string{"!t", "debit2", "debit1"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!t <transfer_from> <transfer_to> <amount> <description>)"),
		},
		{
			it:           "return error when amount cannot be parsed to float64",
			tokenizedMsg: []string{"!t", "debit2", "debit1", "invalid"},
			expectedErr:  errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			res, err := parseTransferRequest(tc.tokenizedMsg)
			assert.Nil(t, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateLength(t *testing.T) {
	testcases := []struct {
		it            string
		tokenizedMsg  []string
		minLength     int
		commandSyntax string
		expectedErr   *errors.AppError
	}{
		{
			it:           "return nil when command's length is valid",
			tokenizedMsg: []string{"!p", "debit1", "200sh"},
			minLength:    3,
			expectedErr:  nil,
		},
		{
			it:            "return error when command's length is less than minLength",
			tokenizedMsg:  []string{"!p", "debit1"},
			minLength:     3,
			commandSyntax: "!p/!e <account_name> <amount><category> <description>",
			expectedErr:   errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax (!p/!e <account_name> <amount><category> <description>)"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			err := validateLength(tc.tokenizedMsg, tc.minLength, tc.commandSyntax)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
