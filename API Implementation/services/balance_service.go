package services

import (
	"CoreImplementation/db"
	"CoreImplementation/models"
	"database/sql"
	"time"
)

type balanceRepo struct {
	database *sql.DB
}

func (r *balanceRepo) Get(userID int) (*models.Balance, error) {
	return db.GetBalance(r.database, userID)
}

func (r *balanceRepo) Update(b *models.Balance) error {
	return db.UpdateBalance(r.database, b)
}

func (r *balanceRepo) History(userID int) ([]models.Balance, error) {
	return models.GetHistory(userID), nil
}

type BalanceService struct {
	repo *balanceRepo
}

func NewBalanceService(database *sql.DB) *BalanceService {
	return &BalanceService{
		repo: &balanceRepo{database: database},
	}
}

func (s *BalanceService) Update(userID int, amt float64) error {
	b, err := s.repo.Get(userID)
	if err != nil {
		return err
	}
	b.Update(amt)
	return s.repo.Update(b)
}

func (s *BalanceService) Get(userID int) (*models.Balance, error) {
	return s.repo.Get(userID)
}

func (s *BalanceService) History(userID int) ([]models.Balance, error) {
	return s.repo.History(userID)
}

func (s *BalanceService) GetCurrentBalance(userID int) (*models.Balance, error) {
	return s.repo.Get(userID)
}

func (s *BalanceService) GetBalanceHistory(userID int) ([]models.Balance, error) {
	return s.repo.History(userID)
}

func (s *BalanceService) GetBalanceAtTime(userID int, t time.Time) (*models.Balance, error) {
	history, err := s.repo.History(userID)
	if err != nil {
		return nil, err
	}
	var closest *models.Balance
	minDiff := time.Duration(1<<63 - 1)
	for _, b := range history {
		diff := t.Sub(b.LastUpdatedAt)
		if diff < 0 {
			diff = -diff
		}
		if diff < minDiff {
			minDiff = diff
			copyB := b
			closest = &copyB
		}
	}
	if closest == nil {
		return nil, nil
	}
	return closest, nil
}
