package repository

import (
	"context"
	"database/sql"

	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/model"
)

type TransactionRepository interface {
	CreateTransactionTx(ctx context.Context, tx *sql.Tx, tr *model.Transaction) (int64, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransactionTx(ctx context.Context, tx *sql.Tx, tr *model.Transaction) (int64, error) {
	query := "INSERT INTO transactions (user_id, amount, type, created_at) VALUES (?, ?, ?, NOW())"
	result, err := tx.ExecContext(ctx, query, tr.UserID, tr.Amount, tr.Type)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
