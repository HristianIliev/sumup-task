package dbmodels

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	Hash       string   `db:"hash"`
	Status     byte     `db:"status"`
	BlocHash   string   `db:"bloc_hash"`
	BlocNumber *big.Int `db:"bloc_number"`
	Sender     string   `db:"sender"`
	Recipient  string   `db:"recipient"`
	Address    string   `db:"address"`
	LogsCount  int      `db:"logs_count"`
	Input      string   `db:"input"`
	Value      string   `db:"value"`
}

func From(r *types.Receipt, t *types.Transaction) (*Transaction, error) {
	from, err := types.Sender(types.LatestSignerForChainID(t.ChainId()), t)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Hash:       r.TxHash.String(),
		Status:     byte(r.Status),
		BlocHash:   r.BlockHash.String(),
		BlocNumber: r.BlockNumber,
		Sender:     from.String(),
		Recipient:  t.To().String(),
		Address:    r.ContractAddress.String(),
		LogsCount:  len(r.Logs),
		Value:      t.Value().String(),
	}, nil
}
