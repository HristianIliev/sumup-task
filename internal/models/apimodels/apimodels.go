package apimodels

import "ethereum-fetcher/internal/models/dbmodels"

type TransactionsResponse struct {
	transaction []*transaction
}

type transaction struct {
}

func From(transactions []*dbmodels.Transaction) *TransactionsResponse {
	return nil
}
