package services

import (
	"CoreImplementation/db"
	"CoreImplementation/models"
	"database/sql"
)

type txRepo struct {
	database *sql.DB
}

func (r *txRepo) Create(tx *models.Transaction) error {
	return db.InsertTransaction(r.database, tx)
}

func (r *txRepo) UpdateStatus(id int, status string) error {
	return db.UpdateTransactionStatus(r.database, id, status)
}

type TransactionService struct {
	repo *txRepo
}

func NewTransactionService(database *sql.DB) *TransactionService {
	return &TransactionService{
		repo: &txRepo{database: database},
	}
}

func (s *TransactionService) Transfer(from, to int, amt float64) (*models.Transaction, error) {
	tx := &models.Transaction{
		FromUserID: from,
		ToUserID:   to,
		Amount:     amt,
		Status:     models.TxPending,
	}
	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *TransactionService) Rollback(txID int) error {
	return s.repo.UpdateStatus(txID, models.TxFailed)
}
