package finance

import (
	"fmt"
	"strings"
	"time"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/internal/ports/client"
)

// Handler implements command handling for finance-related commands.
type Handler struct {
	client client.FinanceServiceClient
}

// NewHandler constructs a finance command handler.
func NewHandler(client client.FinanceServiceClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Match(cmd string) bool {
	if _, exist := commandPrefixSet[cmd]; exist {
		return true
	}
	return false
}

func (h *Handler) Handle(tokenizedMsg []string) (string, *errors.AppError) {
	if len(tokenizedMsg) == 0 {
		return "", errors.BadRequestError(invalidCommandMsg)
	}
	switch tokenizedMsg[0] {
	case "!p":
		return h.withdraw(tokenizedMsg)
	case "!e":
		return h.deposit(tokenizedMsg)
	case "!t":
		return h.transfer(tokenizedMsg)
	case "balance":
		return h.getBalance()
	case "statement":
		return h.getStatement(tokenizedMsg)
	default:
		return "", errors.BadRequestError(invalidCommandMsg)
	}
}

func (h *Handler) withdraw(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransactionRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Withdraw(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully withdraw\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}

func (h *Handler) deposit(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransactionRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Deposit(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully deposit\n================\nResult\nAccount: %v\nBalance: ฿%v", res.Account, res.Balance), nil
}

func (h *Handler) transfer(tokenizedMsg []string) (string, *errors.AppError) {
	req, err := parseTransferRequest(tokenizedMsg)
	if err != nil {
		return "", err
	}
	res, err := h.client.Transfer(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Succesfully transfer\n================\nResult\nAccount: %v\nBalance: ฿%v", res.FromAccount, res.Balance), nil
}

func (h *Handler) getBalance() (string, *errors.AppError) {
	res, err := h.client.GetBalance()
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

func (h *Handler) getStatement(tokenizedMsg []string) (string, *errors.AppError) {
	var res *domain.GetOverviewStatementResponse
	var err *errors.AppError
	statementType := "Income"
	switch len(tokenizedMsg) {
	case 1:
		res, statementType, err = h.callMonthlyOrAnnualStatement("m")
	case 2:
		res, statementType, err = h.callMonthlyOrAnnualStatement(tokenizedMsg[1])
	case 3:
		res, err = h.callSelectedRangeStatement(tokenizedMsg[1], tokenizedMsg[2])
	default:
		err = errors.BadRequestError("Invalid command")
	}

	if err != nil {
		return "", err
	}
	return printStatement(res, statementType)
}

func (h *Handler) callMonthlyOrAnnualStatement(statmentType string) (*domain.GetOverviewStatementResponse, string, *errors.AppError) {
	switch statmentType {
	case "m":
		res, err := h.client.GetOverviewMonthlyStatement()
		return res, "Monthly", err
	case "a":
		res, err := h.client.GetOverviewAnnualStatement()
		return res, "Annual", err
	default:
		return nil, "", errors.BadRequestError(invalidCommandMsg)
	}
}

func (h *Handler) callSelectedRangeStatement(from, to string) (*domain.GetOverviewStatementResponse, *errors.AppError) {
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
	return h.client.GetOverviewStatement(req)
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
