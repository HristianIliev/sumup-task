package controllers

import (
	"context"
	"ethereum-fetcher/internal/models/apimodels"
	"ethereum-fetcher/internal/models/dbmodels"
	"ethereum-fetcher/internal/service/transaction"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionController struct {
	transactionService transaction.TransactionService
}

func New(transactionService transaction.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// Subsequent calls for this transaction identified by its hash should get the information from the database rather than the Ethereum node.
// Use LRU cache for that
func (t *TransactionController) GetTransactions(hashes []string) (*apimodels.TransactionsResponse, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/71680f2f4b574e338a2a65fea4340b1b")
	if err != nil {
		log.Fatal(err)
	}

	transactions := []*dbmodels.Transaction{}
	for _, hash := range hashes {
		receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(hash))
		if err != nil {
			log.Printf("%+v", err)
		}

		transaction, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(hash))
		if err != nil || isPending {
			log.Printf("%+v", err)
		}

		if transaction == nil || receipt == nil {
			continue
		}

		if t, err := dbmodels.From(receipt, transaction); err == nil {
			transactions = append(transactions, t)
		}
	}

	if err := t.transactionService.Store(transactions); err != nil {
		return nil, err
	}

	return apimodels.From(transactions), nil
}
