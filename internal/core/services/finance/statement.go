package finance

import (
	"fmt"
	"strings"
	"time"

	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/secretaria-bot/internal/core/errors"
)

// TODO: Refactor
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
	return printStatement(res, statementType), nil
}

// TODO: Refactor
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

// TODO: Refactor time in the database to be in UTC
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

func printStatement(res *domain.GetOverviewStatementResponse, statementType string) string {
	if res.Revenue == nil {
		res.Revenue = &domain.GetOverviewStatementSection{}
	}
	if res.Expense == nil {
		res.Expense = &domain.GetOverviewStatementSection{}
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
	return sb.String()
}
