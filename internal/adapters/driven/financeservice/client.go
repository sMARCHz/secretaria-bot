package financeservice

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driven/financeservice/pb"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/dto"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type financeServiceClient struct {
	client pb.FinanceServiceClient
	logger logger.Logger
}

func NewFinanceServiceClient(url string, logger logger.Logger) client.FinanceServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("could not connect to %v:", url, err)
	}
	return &financeServiceClient{
		client: pb.NewFinanceServiceClient(conn),
		logger: logger,
	}
}

func (f *financeServiceClient) Withdraw(msg []string) (*dto.TransactionResponse, *errors.AppError) {
	req, appErr := f.newTransactionRequest(msg)
	if appErr != nil {
		return nil, appErr
	}
	res, err := f.client.Withdraw(context.Background(), req)
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return &dto.TransactionResponse{
		AccountName: res.AccountName,
		Balance:     res.Balance,
	}, nil
}

func (f *financeServiceClient) Deposit(msg []string) (*dto.TransactionResponse, *errors.AppError) {
	req, appErr := f.newTransactionRequest(msg)
	if appErr != nil {
		return nil, appErr
	}
	res, err := f.client.Deposit(context.Background(), req)
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return &dto.TransactionResponse{
		AccountName: res.AccountName,
		Balance:     res.Balance,
	}, nil
}

func (f *financeServiceClient) Transfer(msg []string) (*dto.TransferResponse, *errors.AppError) {
	req, appErr := f.newTransferRequest(msg)
	if appErr != nil {
		return nil, appErr
	}
	res, err := f.client.Transfer(context.Background(), req)
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return &dto.TransferResponse{
		FromAccountName: res.FromAccountName,
		Balance:         res.Balance,
	}, nil
}

func (f *financeServiceClient) GetBalance() (*dto.GetBalanceResponse, *errors.AppError) {
	res, err := f.client.GetBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	accounts := make([]dto.AccountBalance, len(res.Accounts))
	for i, v := range res.Accounts {
		accounts[i] = dto.AccountBalance{AccountName: v.AccountName, Balance: v.Balance}
	}
	return &dto.GetBalanceResponse{
		Accounts: accounts,
	}, nil
}

func (f *financeServiceClient) GetOverviewStatement(from time.Time, to time.Time) (*dto.GetOverviewStatementResponse, *errors.AppError) {
	req := &pb.OverviewStatementRequest{
		From: timestamppb.New(from),
		To:   timestamppb.New(to),
	}
	res, err := f.client.GetOverviewStatement(context.Background(), req)
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponseDto(res), nil
}

func (f *financeServiceClient) GetOverviewMonthlyStatement() (*dto.GetOverviewStatementResponse, *errors.AppError) {
	res, err := f.client.GetOverviewMonthlyStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponseDto(res), nil
}

func (f *financeServiceClient) GetOverviewAnnualStatement() (*dto.GetOverviewStatementResponse, *errors.AppError) {
	res, err := f.client.GetOverviewAnnualStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponseDto(res), nil
}

func (f *financeServiceClient) newTransactionRequest(msg []string) (*pb.TransactionRequest, *errors.AppError) {
	size := len(msg)
	if size < 3 {
		f.logger.Error("Invalid command length")
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax, !p/!e <account_name> <amount><category> <description>")
	}

	regex := regexp.MustCompile(`(^\d+(\.\d+)?)([a-z]+)`)
	amountAndCategory := regex.FindStringSubmatch(msg[2])
	var amountAsStr string
	var category string
	acSize := len(amountAndCategory)
	if acSize == 3 {
		amountAsStr = amountAndCategory[1]
		category = amountAndCategory[2]
	} else if acSize == 4 {
		amountAsStr = amountAndCategory[1]
		category = amountAndCategory[3]
	} else {
		f.logger.Error("invalid amount and category combination['%v']", msg[2])
		return nil, errors.BadRequestError("Invalid amount and category combination")
	}

	amount, err := strconv.ParseFloat(amountAsStr, 64)
	if err != nil {
		f.logger.Error("cannot parse amount to float64: ", err)
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command")
	}

	var description string
	if size > 3 {
		description = strings.Join(msg[3:], " ")
	}

	req := &pb.TransactionRequest{
		AccountName: msg[1],
		Amount:      amount,
		Category:    category,
		Description: description,
	}
	return req, nil
}

func (f *financeServiceClient) newTransferRequest(msg []string) (*pb.TransferRequest, *errors.AppError) {
	size := len(msg)
	if size < 4 {
		f.logger.Error("invalid command length")
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck the syntax, !t <transfer_from> <transfer_to> <amount> <description>")
	}

	amount, err := strconv.ParseFloat(msg[3], 64)
	if err != nil {
		f.logger.Error("cannot parse amount to float64: ", err)
		return nil, errors.BadRequestError("Invalid command's arguments.\nPlease recheck syntax and amount of transaction in the command")
	}

	var description string
	if size > 4 {
		description = strings.Join(msg[4:], " ")
	}

	req := &pb.TransferRequest{
		FromAccountName: msg[1],
		ToAccountName:   msg[2],
		Amount:          amount,
		Description:     description,
	}
	return req, nil
}

func (*financeServiceClient) toGetOverviewStatementResponseDto(o *pb.OverviewStatementResponse) *dto.GetOverviewStatementResponse {
	revenueEntries := make([]dto.CategorizedEntry, len(o.Revenue.Entries))
	expenseEntries := make([]dto.CategorizedEntry, len(o.Expense.Entries))
	for i, v := range o.Revenue.Entries {
		revenueEntries[i] = dto.CategorizedEntry{Category: v.Category, Amount: v.Amount}
	}
	for i, v := range o.Expense.Entries {
		expenseEntries[i] = dto.CategorizedEntry{Category: v.Category, Amount: v.Amount}
	}
	revenue := dto.GetOverviewStatementSection{
		Total:   o.Revenue.Total,
		Entries: revenueEntries,
	}
	expense := dto.GetOverviewStatementSection{
		Total:   o.Expense.Total,
		Entries: expenseEntries,
	}
	return &dto.GetOverviewStatementResponse{
		Revenue: revenue,
		Expense: expense,
		Profit:  o.Profit,
	}
}
