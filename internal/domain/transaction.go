package domain

import "context"

// TransactionManager handles database transactions.
type TransactionManager interface {
	// Do executes the given function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
