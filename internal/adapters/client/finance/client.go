package finance

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sMARCHz/secretaria-bot/internal/adapters/client/finance/pb"
	"github.com/sMARCHz/secretaria-bot/internal/config"
	"github.com/sMARCHz/secretaria-bot/internal/core/domain"
	apperrors "github.com/sMARCHz/secretaria-bot/internal/core/errors"
	"github.com/sMARCHz/secretaria-bot/internal/logger"
	"github.com/sMARCHz/secretaria-bot/internal/ports/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type financeServiceClient struct {
	client pb.FinanceServiceClient
}

func NewFinanceServiceClient() client.FinanceServiceClient {
	url := config.Get().FinanceServiceURL
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("could not connect to %v: %v", url, err)
	}
	return &financeServiceClient{
		client: pb.NewFinanceServiceClient(conn),
	}
}

func (f *financeServiceClient) Withdraw(req *domain.TransactionRequest) (*domain.TransactionResponse, *apperrors.AppError) {
	res, err := f.client.Withdraw(context.Background(), req.ToProto())
	if err != nil {
		err = errors.Wrap(err, "cannot withdraw money")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return &domain.TransactionResponse{
		Account: res.AccountName,
		Balance: res.Balance,
	}, nil
}

func (f *financeServiceClient) Deposit(req *domain.TransactionRequest) (*domain.TransactionResponse, *apperrors.AppError) {
	res, err := f.client.Deposit(context.Background(), req.ToProto())
	if err != nil {
		err = errors.Wrap(err, "cannot deposit money")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return &domain.TransactionResponse{
		Account: res.AccountName,
		Balance: res.Balance,
	}, nil
}

func (f *financeServiceClient) Transfer(req *domain.TransferRequest) (*domain.TransferResponse, *apperrors.AppError) {
	res, err := f.client.Transfer(context.Background(), req.ToProto())
	if err != nil {
		err = errors.Wrap(err, "cannot transfer money")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return &domain.TransferResponse{
		FromAccount: res.FromAccountName,
		Balance:     res.Balance,
	}, nil
}

func (f *financeServiceClient) GetBalance() (*domain.GetBalanceResponse, *apperrors.AppError) {
	res, err := f.client.GetBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		err = errors.Wrap(err, "cannot get balance")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	accounts := make([]domain.AccountBalance, len(res.Accounts))
	for i, v := range res.Accounts {
		accounts[i] = domain.AccountBalance{Account: v.AccountName, Balance: v.Balance}
	}
	return &domain.GetBalanceResponse{
		Accounts: accounts,
	}, nil
}

func (f *financeServiceClient) GetOverviewStatement(req *domain.GetOverviewStatementRequest) (*domain.GetOverviewStatementResponse, *apperrors.AppError) {
	res, err := f.client.GetOverviewStatement(context.Background(), req.ToProto())
	if err != nil {
		err = errors.Wrap(err, "cannot get overview statement")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (f *financeServiceClient) GetOverviewMonthlyStatement() (*domain.GetOverviewStatementResponse, *apperrors.AppError) {
	res, err := f.client.GetOverviewMonthlyStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		err = errors.Wrap(err, "cannot get monthly overview statement")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (f *financeServiceClient) GetOverviewAnnualStatement() (*domain.GetOverviewStatementResponse, *apperrors.AppError) {
	res, err := f.client.GetOverviewAnnualStatement(context.Background(), &emptypb.Empty{})
	if err != nil {
		err = errors.Wrap(err, "cannot get annual overview statement")
		logger.Error(err)
		return nil, apperrors.BadGatewayError(err.Error())
	}
	return f.toGetOverviewStatementResponse(res), nil
}

func (*financeServiceClient) toGetOverviewStatementResponse(o *pb.OverviewStatementResponse) *domain.GetOverviewStatementResponse {
	if o == nil {
		return &domain.GetOverviewStatementResponse{}
	}

	var revenue *domain.GetOverviewStatementSection
	if o.Revenue != nil {
		entries := make([]domain.CategorizedEntry, len(o.Revenue.Entries))
		for i, v := range o.Revenue.Entries {
			entries[i] = domain.CategorizedEntry{Category: v.Category, Amount: v.Amount}
		}
		revenue = &domain.GetOverviewStatementSection{
			Entries: entries,
			Total:   o.Revenue.GetTotal(),
		}
	}

	var expense *domain.GetOverviewStatementSection
	if o.Expense != nil {
		entries := make([]domain.CategorizedEntry, len(o.Expense.Entries))
		for i, v := range o.Expense.Entries {
			entries[i] = domain.CategorizedEntry{Category: v.Category, Amount: v.Amount}
		}
		expense = &domain.GetOverviewStatementSection{
			Entries: entries,
			Total:   o.Expense.GetTotal(),
		}
	}
	return &domain.GetOverviewStatementResponse{
		Revenue: revenue,
		Expense: expense,
		Profit:  o.Profit,
	}
}
