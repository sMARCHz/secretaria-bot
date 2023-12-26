package financeservice

import (
	"context"

	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driven/financeservice/pb"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/client"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/domain"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type financeServiceClient struct {
	client pb.FinanceServiceClient
}

func NewFinanceServiceClient(url string) client.FinanceServiceClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("could not connect to %v:", url, err)
	}
	return &financeServiceClient{
		client: pb.NewFinanceServiceClient(conn),
	}
}

func (f *financeServiceClient) Withdraw(req *domain.TransactionRequest) (*domain.TransactionResponse, *errors.AppError) {
	res, err := f.client.Withdraw(context.Background(), req.ToProto())
	if err != nil {
		logger.Error("cannot withdraw money: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return &domain.TransactionResponse{
		Account: res.AccountName,
		Balance: res.Balance,
	}, nil
}

func (f *financeServiceClient) Deposit(req *domain.TransactionRequest) (*domain.TransactionResponse, *errors.AppError) {
	res, err := f.client.Deposit(context.Background(), req.ToProto())
	if err != nil {
		logger.Error("cannot deposit money: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return &domain.TransactionResponse{
		Account: res.AccountName,
		Balance: res.Balance,
	}, nil
}

func (f *financeServiceClient) Transfer(req *domain.TransferRequest) (*domain.TransferResponse, *errors.AppError) {
	res, err := f.client.Transfer(context.Background(), req.ToProto())
	if err != nil {
		logger.Error("cannot transfer money: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return &domain.TransferResponse{
		FromAccount: res.FromAccountName,
		Balance:     res.Balance,
	}, nil
}

func (f *financeServiceClient) GetBalance() (*domain.GetBalanceResponse, *errors.AppError) {
	res, err := f.client.GetBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Error("cannot get balance: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	accounts := make([]domain.AccountBalance, len(res.Accounts))
	for i, v := range res.Accounts {
		accounts[i] = domain.AccountBalance{Account: v.AccountName, Balance: v.Balance}
	}
	return &domain.GetBalanceResponse{
		Accounts: accounts,
	}, nil
}

func (f *financeServiceClient) GetOverviewStatement(req *domain.GetOverviewStatementRequest) (*domain.GetOverviewStatementResponse, *errors.AppError) {
	res, err := f.client.GetOverviewStatement(context.Background(), req.ToProto())
	if err != nil {
		logger.Error("cannot get overview statement: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (f *financeServiceClient) GetOverviewMonthlyStatement() (*domain.GetOverviewStatementResponse, *errors.AppError) {
	res, err := f.client.GetOverviewMonthlyStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Error("cannot get monthly overview statement: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (f *financeServiceClient) GetOverviewAnnualStatement() (*domain.GetOverviewStatementResponse, *errors.AppError) {
	res, err := f.client.GetOverviewAnnualStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Error("cannot get annual overview statement: ", err)
		return nil, errors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (*financeServiceClient) toGetOverviewStatementResponse(o *pb.OverviewStatementResponse) *domain.GetOverviewStatementResponse {
	revenueEntries := make([]domain.CategorizedEntry, len(o.Revenue.Entries))
	expenseEntries := make([]domain.CategorizedEntry, len(o.Expense.Entries))
	for i, v := range o.Revenue.Entries {
		revenueEntries[i] = domain.CategorizedEntry{Category: v.Category, Amount: v.Amount}
	}
	for i, v := range o.Expense.Entries {
		expenseEntries[i] = domain.CategorizedEntry{Category: v.Category, Amount: v.Amount}
	}
	return &domain.GetOverviewStatementResponse{
		Revenue: domain.GetOverviewStatementSection{
			Total:   o.Revenue.Total,
			Entries: revenueEntries,
		},
		Expense: domain.GetOverviewStatementSection{
			Total:   o.Expense.Total,
			Entries: expenseEntries,
		},
		Profit: o.Profit,
	}
}
