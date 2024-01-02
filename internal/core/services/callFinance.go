package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/utils"
)

func (b *botService) callWithdraw(msgArr []string) (string, *errors.AppError) {
	req, err := utils.ParseTransactionRequest(msgArr)
	if err != nil {
		return "", err
	}
	res, err := b.financeService.Withdraw(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully withdraw\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}

func (b *botService) callDeposit(msgArr []string) (string, *errors.AppError) {
	req, err := utils.ParseTransactionRequest(msgArr)
	if err != nil {
		return "", err
	}
	res, err := b.financeService.Deposit(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully deposit\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}

func (b *botService) callTransfer(msgArr []string) (string, *errors.AppError) {
	req, err := utils.ParseTransferRequest(msgArr)
	if err != nil {
		return "", err
	}
	res, err := b.financeService.Transfer(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully transfer\n================\nResult\nAccount: %v\nBalance: ฿%v", res.FromAccount, res.Balance), nil
}

func (b *botService) callBalance() (string, *errors.AppError) {
	res, err := b.financeService.GetBalance()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Your balance\n\n")
	for _, v := range res.Accounts {
		sb.WriteString(fmt.Sprintf("Account: %v => Balance: ฿%v\n", v.Account, v.Balance))
	}
	return sb.String(), nil
}

func (b *botService) callStatement(msgArr []string) (string, *errors.AppError) {
	var res *domain.GetOverviewStatementResponse
	var err *errors.AppError
	statementType := "Income"
	switch len(msgArr) {
	case 1:
		res, statementType, err = b.callMonthlyOrAnnualStatement("m")
	case 2:
		res, statementType, err = b.callMonthlyOrAnnualStatement(msgArr[1])
	case 3:
		res, err = b.callSelectedRangeStatement(msgArr[1], msgArr[2])
	default:
		err = errors.BadRequestError("Invalid command")
	}

	if err != nil {
		return "", err
	}
	return printStatement(res, statementType)
}

func (b *botService) callMonthlyOrAnnualStatement(statmentType string) (*domain.GetOverviewStatementResponse, string, *errors.AppError) {
	switch statmentType {
	case "m":
		res, err := b.financeService.GetOverviewMonthlyStatement()
		return res, "Monthly", err
	case "a":
		res, err := b.financeService.GetOverviewAnnualStatement()
		return res, "Annual", err
	default:
		return nil, "", errors.BadRequestError("Invalid command")
	}
}

func (b *botService) callSelectedRangeStatement(from, to string) (*domain.GetOverviewStatementResponse, *errors.AppError) {
	fromAsTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the from_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>")
	}
	toAsTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the to_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>")
	}
	req := &domain.GetOverviewStatementRequest{
		From: fromAsTime,
		To:   toAsTime,
	}
	return b.financeService.GetOverviewStatement(req)
}

func printStatement(res *domain.GetOverviewStatementResponse, statementType string) (string, *errors.AppError) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v Statement\n================\n", statementType))
	sb.WriteString(fmt.Sprintf("Revenue: ฿%v\n", res.Revenue.Total))
	for _, v := range res.Revenue.Entries {
		sb.WriteString(fmt.Sprintf("%v = ฿%v\n", v.Category, v.Amount))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Expense: ฿%v\n", res.Expense.Total))
	for _, v := range res.Expense.Entries {
		sb.WriteString(fmt.Sprintf("%v = ฿%v\n", v.Category, v.Amount))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Profit: ฿%v", res.Profit))
	return sb.String(), nil
}
