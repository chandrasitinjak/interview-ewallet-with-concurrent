package service

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/model"
	"github.com/chandrasitinjak/interview-ewallet-with-concurrent/repository"
)

type TransactionService interface {
	Credit(ctx context.Context, userID int64, amount float64) (*model.Transaction, float64, error)
	Debit(ctx context.Context, userID int64, amount float64) (*model.Transaction, float64, error)
}

type transactionService struct {
	userRepo repository.UserRepository
	txRepo   repository.TransactionRepository
	mu       sync.Mutex
}

func NewTransactionService(userRepo repository.UserRepository, txRepo repository.TransactionRepository) TransactionService {
	return &transactionService{userRepo: userRepo, txRepo: txRepo}
}

func (s *transactionService) Credit(ctx context.Context, userID int64, amount float64) (*model.Transaction, float64, error) {
	if amount <= 0 {
		return nil, 0, errors.New("invalid amount")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Begin transaction
	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback() // rollback jika gagal

	// Lock user row
	_, err = s.userRepo.GetUserForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Update balance
	if err := s.userRepo.UpdateBalanceTx(ctx, tx, userID, amount); err != nil {
		return nil, 0, err
	}

	// Create transaction
	trx := &model.Transaction{
		UserID: userID,
		Amount: amount,
		Type:   "credit",
	}
	id, err := s.txRepo.CreateTransactionTx(ctx, tx, trx)
	if err != nil {
		return nil, 0, err
	}
	trx.ID = id

	// Ambil balance terbaru
	newUser, err := s.userRepo.GetUserForUpdate(ctx, tx, userID)
	if err != nil {
		log.Print("error123 : ", err.Error())
		return &model.Transaction{}, 0, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	return trx, newUser.Balance, nil
}

func (s *transactionService) Debit(ctx context.Context, userID int64, amount float64) (*model.Transaction, float64, error) {
	if amount <= 0 {
		return nil, 0, errors.New("invalid amount")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Begin transaction
	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback() // rollback otomatis kalau gagal

	// Lock row user
	user, err := s.userRepo.GetUserForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Check saldo cukup
	if user.Balance < amount {
		return nil, user.Balance, errors.New("insufficient funds")
	}

	// Update saldo (dikurangi)
	if err := s.userRepo.UpdateBalanceTx(ctx, tx, userID, -amount); err != nil {
		return nil, 0, err
	}

	// Catat transaksi
	trx := &model.Transaction{
		UserID: userID,
		Amount: amount,
		Type:   "debit",
	}
	id, err := s.txRepo.CreateTransactionTx(ctx, tx, trx)
	if err != nil {
		return nil, 0, err
	}
	trx.ID = id

	// Ambil saldo terbaru
	newUser, err := s.userRepo.GetUserForUpdate(ctx, tx, userID)
	if err != nil {
		log.Print("error123 : ", err.Error())
		return &model.Transaction{}, 0, err
	}

	// Commit perubahan
	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	return trx, newUser.Balance, nil
}
