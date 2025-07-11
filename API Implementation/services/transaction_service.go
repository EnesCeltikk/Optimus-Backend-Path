package services

import (
	"CoreImplementation/db"
	"CoreImplementation/models"
	"database/sql"
	"time"
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

func (s *TransactionService) Credit(toUserID int, amount float64) error {
	tx := &models.Transaction{
		FromUserID: 0,
		ToUserID:   toUserID,
		Amount:     amount,
		Type:       "credit",
		Status:     models.TxSuccess,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(tx)
}

func (s *TransactionService) Debit(fromUserID int, amount float64) error {
	tx := &models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   0,
		Amount:     amount,
		Type:       "debit",
		Status:     models.TxSuccess,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(tx)
}

func (s *TransactionService) GetTransactionHistory(userID, page, limit int) ([]*models.Transaction, error) {
	offset := (page - 1) * limit
	if page < 1 {
		offset = 0
	}
	rows, err := s.repo.database.Query("SELECT id, from_user_id, to_user_id, amount, type, status, created_at FROM transactions WHERE from_user_id=$1 OR to_user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txs []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (s *TransactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	row := s.repo.database.QueryRow("SELECT id, from_user_id, to_user_id, amount, type, status, created_at FROM transactions WHERE id=$1", id)
	tx := &models.Transaction{}
	if err := row.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
		return nil, err
	}
	return tx, nil
}
