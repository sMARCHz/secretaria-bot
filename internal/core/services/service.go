package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/core/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/utils"
)

type BotService interface {
	HandleTextMessage(string) (*domain.TextMessageResponse, *errors.AppError)
}

type botService struct {
	financeService client.FinanceServiceClient
}

func NewBotService(financeService client.FinanceServiceClient) BotService {
	return &botService{
		financeService: financeService,
	}
}

func (b *botService) HandleTextMessage(msg string) (*domain.TextMessageResponse, *errors.AppError) {
	replyMsg := ""
	msg = strings.TrimSpace(msg)
	msg = strings.ToLower(msg)
	msgArr := strings.Fields(msg)

	var err *errors.AppError
	switch msgArr[0] {
	case "!p":
		replyMsg, err = b.callWithdraw(msgArr)
	case "!e":
		replyMsg, err = b.callDeposit(msgArr)
	case "!t":
		replyMsg, err = b.callTransfer(msgArr)
	case "balance":
		replyMsg, err = b.callBalance()
	case "statement":
		replyMsg, err = b.callStatement(msgArr)
	default:
		replyMsg = "Command not found"
	}

	if err != nil {
		return nil, err
	}
	return &domain.TextMessageResponse{ReplyMessage: replyMsg}, nil
}

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
	var appErr *errors.AppError
	var statementType string
	if len(msgArr) == 1 || msgArr[1] == "m" {
		statementType = "Monthly"
		res, appErr = b.financeService.GetOverviewMonthlyStatement()
	} else if msgArr[1] == "a" {
		statementType = "Annual"
		res, appErr = b.financeService.GetOverviewAnnualStatement()
	} else {
		statementType = "Income"
		from, err := time.Parse("2006-01-02", msgArr[1])
		if err != nil {
			return "", errors.BadRequestError("Invalid command's arguments.\nPlease recheck the from_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>")
		}
		to, err := time.Parse("2006-01-02", msgArr[2])
		if err != nil {
			return "", errors.BadRequestError("Invalid command's arguments.\nPlease recheck the to_date, <statement> <from_date: 2022-01-01> <to_date: 2022-01-01>")
		}
		req := &domain.GetOverviewStatementRequest{
			From: from,
			To:   to,
		}
		res, appErr = b.financeService.GetOverviewStatement(req)
	}
	if appErr != nil {
		return "", appErr
	}

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
