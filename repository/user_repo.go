package repository

import (
	"context"
	"database/sql"

	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/model"
)

type UserRepository interface {
	GetUserForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error)
	UpdateBalanceTx(ctx context.Context, tx *sql.Tx, userID int64, amount float64) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *userRepository) GetUserForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error) {
	query := "SELECT id, username, balance, created_at FROM users WHERE id = ? FOR UPDATE"
	row := tx.QueryRowContext(ctx, query, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Username, &u.Balance, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) UpdateBalanceTx(ctx context.Context, tx *sql.Tx, userID int64, amount float64) error {
	query := "UPDATE users SET balance = balance + ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, amount, userID)
	return err
}
