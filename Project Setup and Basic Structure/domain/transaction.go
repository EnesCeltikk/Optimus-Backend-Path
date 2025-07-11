package domain

import (
	"errors"
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRolledBack TransactionStatus = "rolled_back"
)

type Transaction struct {
	ID         int               `json:"id"`
	FromUserID int               `json:"from_user_id"`
	ToUserID   int               `json:"to_user_id"`
	Amount     float64           `json:"amount"`
	Type       TransactionType   `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if t.Type == TransactionTypeTransfer && (t.FromUserID == 0 || t.ToUserID == 0) {
		return errors.New("both from_user_id and to_user_id are required for transfers")
	}
	return nil
}

func (t *Transaction) CanTransitionTo(newStatus TransactionStatus) bool {
	switch t.Status {
	case TransactionStatusPending:
		return newStatus == TransactionStatusCompleted || newStatus == TransactionStatusFailed
	case TransactionStatusCompleted:
		return newStatus == TransactionStatusRolledBack
	default:
		return false
	}
}

func (t *Transaction) TransitionTo(newStatus TransactionStatus) error {
	if !t.CanTransitionTo(newStatus) {
		return errors.New("invalid status transition")
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

func (t *Transaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

func (t *Transaction) IsRolledBack() bool {
	return t.Status == TransactionStatusRolledBack
} 