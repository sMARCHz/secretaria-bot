package dto

// Transaction
type TransactionResponse struct {
	AccountName string
	Balance     float64
}

// Transfer
type TransferResponse struct {
	FromAccountName string
	Balance         float64
}

// GetBalance
type GetBalanceResponse struct {
	Accounts []AccountBalance
}

type AccountBalance struct {
	AccountName string
	Balance     float64
}

// GetOverviewStatementResponse
type GetOverviewStatementResponse struct {
	Revenue GetOverviewStatementSection
	Expense GetOverviewStatementSection
	Profit  float64
}

type GetOverviewStatementSection struct {
	Total   float64
	Entries []CategorizedEntry
}

type CategorizedEntry struct {
	Category string
	Amount   float64
}
