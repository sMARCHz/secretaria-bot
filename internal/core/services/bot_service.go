package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/dto"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

type BotService interface {
	HandleTextMessage(string) (*dto.TextMessageResponse, *errors.AppError)
}

type botService struct {
	financeService client.FinanceServiceClient
	config         config.Configuration
	logger         logger.Logger
}

func NewBotService(financeService client.FinanceServiceClient, config config.Configuration, logger logger.Logger) BotService {
	return &botService{
		financeService: financeService,
		config:         config,
		logger:         logger,
	}
}

func (b *botService) HandleTextMessage(msg string) (*dto.TextMessageResponse, *errors.AppError) {
	replyMsg := ""
	msg = strings.TrimSpace(msg)
	msg = strings.ToLower(msg)
	msgArr := strings.Fields(msg)
	switch msgArr[0] {
	case "!p", "!e":
		var transactionType string
		var res *dto.TransactionResponse
		var err *errors.AppError
		if msgArr[0] == "!p" {
			transactionType = "withdraw"
			res, err = b.financeService.Withdraw(msgArr)
		} else {
			transactionType = "deposit"
			res, err = b.financeService.Deposit(msgArr)
		}
		if err != nil {
			return nil, err
		}
		replyMsg = fmt.Sprintf("Succesfully %v\n================\nResult\nAccount: %v\nBalance: ฿%v", transactionType, res.AccountName, res.Balance)

	case "!t":
		res, err := b.financeService.Transfer(msgArr)
		if err != nil {
			return nil, err
		}
		replyMsg = fmt.Sprintf("Succesfully transfer\n================\nResult\nAccount: %v\nBalance: ฿%v", res.FromAccountName, res.Balance)

	case "balance":
		res, err := b.financeService.GetBalance()
		if err != nil {
			return nil, err
		}
		var sb strings.Builder
		sb.WriteString("Your balance\n\n")
		for _, v := range res.Accounts {
			sb.WriteString(fmt.Sprintf("Account: %v => Balance: ฿%v\n", v.AccountName, v.Balance))
		}
		replyMsg = sb.String()

	case "statement":
		var res *dto.GetOverviewStatementResponse
		var appErr *errors.AppError
		statementType := "Income"
		if len(msgArr) == 1 || msgArr[1] == "m" {
			statementType = "Monthly"
			res, appErr = b.financeService.GetOverviewMonthlyStatement()
		} else if msgArr[1] == "a" {
			statementType = "Annual"
			res, appErr = b.financeService.GetOverviewAnnualStatement()
		} else {
			from, err := time.Parse("2006-01-02", msgArr[1])
			if err != nil {
				return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax and from_date")
			}
			to, err := time.Parse("2006-01-02", msgArr[2])
			if err != nil {
				return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax and to_date")
			}
			res, appErr = b.financeService.GetOverviewStatement(from, to)
		}
		if appErr != nil {
			return nil, appErr
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
		replyMsg = sb.String()

	default:
		replyMsg = "Command not found"
	}
	return &dto.TextMessageResponse{ReplyMessage: replyMsg}, nil
}
