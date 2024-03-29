package transaction

import (
	"ethereum-fetcher/internal/models/dbmodels"
	"log"

	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

func (t *TransactionService) RetrieveByHashes(transactionHashes []string) ([]*dbmodels.Transaction, error) {
	// TODO: is this injection protected?
	query, args, err := sqlx.In("SELECT * FROM transaction WHERE hash IN (?)", transactionHashes)
	if err != nil {
		log.Fatal(err)
	}

	query = t.db.Rebind(query)

	transactions := make([]*dbmodels.Transaction, 0)
	err = t.db.Select(&transactions, query, args...)
	if err != nil {
		log.Fatal(err)
	}

	return transactions, nil
}

func (t *TransactionService) Store(transactions []*dbmodels.Transaction) error {
	tx := t.db.MustBegin()
	for _, transaction := range transactions {
		tx.NamedExec("INSERT INTO transaction (hash, status, bloc_hash, bloc_number, from, to, address, logs_count, input, value) VALUES (:hash, :status, :bloc_hash, :bloc_number, :from, :to, :address, :logs_count, :input, :value)", transaction)
	}

	return tx.Commit()
}
