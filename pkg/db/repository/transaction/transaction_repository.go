package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
