package domain

import (
	"sync"
	"time"
)

type Balance struct {
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	mu            sync.RWMutex
}

func NewBalance(userID int) *Balance {
	return &Balance{
		UserID:        userID,
		Amount:        0,
		LastUpdatedAt: time.Now(),
	}
}

func (b *Balance) GetAmount() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Amount
}

func (b *Balance) Add(amount float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Amount += amount
	b.LastUpdatedAt = time.Now()
}

func (b *Balance) Subtract(amount float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Amount < amount {
		return false
	}
	b.Amount -= amount
	b.LastUpdatedAt = time.Now()
	return true
}

func (b *Balance) Transfer(to *Balance, amount float64) bool {
	if amount <= 0 {
		return false
	}

	// Lock both balances to prevent deadlock
	if b.UserID < to.UserID {
		b.mu.Lock()
		to.mu.Lock()
	} else {
		to.mu.Lock()
		b.mu.Lock()
	}
	defer b.mu.Unlock()
	defer to.mu.Unlock()

	if b.Amount < amount {
		return false
	}

	b.Amount -= amount
	to.Amount += amount
	now := time.Now()
	b.LastUpdatedAt = now
	to.LastUpdatedAt = now
	return true
}

func (b *Balance) Snapshot() Balance {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return Balance{
		UserID:        b.UserID,
		Amount:        b.Amount,
		LastUpdatedAt: b.LastUpdatedAt,
	}
} 