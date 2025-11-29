package finance

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/internal/logger"
)

var transactionCommandPattern = regexp.MustCompile(`^(\d+(?:\.\d+)?)?([a-zA-Z]+)$`)

// TODO: Rename variable
func parseTransactionRequest(tokenizedMsg []string) (*domain.TransactionRequest, *errors.AppError) {
	if err := validateLength(tokenizedMsg, 3, "!p/!e <account_name> <amount><category> <description>"); err != nil {
		return nil, err
	}

	// 200.12sh -> [200.12sh, 200.12, sh]
	submatch := transactionCommandPattern.FindStringSubmatch(tokenizedMsg[2])
	if len(submatch) != 3 {
		logger.Error("invalid amount and category combination['%v']", tokenizedMsg[2])
		return nil, errors.BadRequestError("Invalid amount and category combination")
	}

	amount, err := strconv.ParseFloat(submatch[1], 64)
	if err != nil {
		logger.Error("cannot parse amount to float64: ", err)
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command")
	}

	var description string
	if len(tokenizedMsg) > 3 {
		description = strings.Join(tokenizedMsg[3:], " ")
	}

	return &domain.TransactionRequest{
		Account:     tokenizedMsg[1],
		Amount:      amount,
		Category:    submatch[2],
		Description: description,
	}, nil
}

func parseTransferRequest(tokenizedMsg []string) (*domain.TransferRequest, *errors.AppError) {
	if err := validateLength(tokenizedMsg, 4, "!t <transfer_from> <transfer_to> <amount> <description>"); err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(tokenizedMsg[3], 64)
	if err != nil {
		logger.Error("cannot parse amount to float64: ", err)
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command")
	}

	var description string
	if len(tokenizedMsg) > 4 {
		description = strings.Join(tokenizedMsg[4:], " ")
	}

	return &domain.TransferRequest{
		FromAccount: tokenizedMsg[1],
		ToAccount:   tokenizedMsg[2],
		Amount:      amount,
		Description: description,
	}, nil
}

func validateLength(tokenizedMsg []string, minLength int, commandSyntax string) *errors.AppError {
	if len(tokenizedMsg) < minLength {
		logger.Error("invalid command length")
		return errors.BadRequestError(fmt.Sprintf("Invalid command's arguments.\nPlease recheck the syntax (%s)", commandSyntax))
	}
	return nil
}
