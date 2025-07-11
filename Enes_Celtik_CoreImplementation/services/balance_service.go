package services

import (
	"CoreImplementation/db"
	"CoreImplementation/models"
	"database/sql"
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
