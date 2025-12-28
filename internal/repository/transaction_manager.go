package repository

import (
	"context"
	"synthetica/internal/domain"

	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &transactionManager{db: db}
}

type txnKey struct{}

// Do executes the given function within a transaction.
func (tm *transactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		// Inject the transaction db into the context
		ctxWithTx := context.WithValue(ctx, txnKey{}, tx)
		return fn(ctxWithTx)
	})
}

// getDB extracts the transaction-aware DB from the context,
// or returns the default DB if no transaction is present.
func getDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(txnKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return defaultDB
}
