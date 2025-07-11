package services

import (
	"CoreImplementation/db"
	"CoreImplementation/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	database *sql.DB
}

func (r *userRepo) Insert(u *models.User) error {
	return db.InsertUser(r.database, u)
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	return db.GetUserByEmail(r.database, email)
}

type UserService struct {
	repo *userRepo
}

func NewUserService(database *sql.DB) *UserService {
	return &UserService{
		repo: &userRepo{database: database},
	}
}

func (s *UserService) Register(username, email, password string) (*models.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := s.repo.Insert(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *UserService) Authorize(u *models.User, action string) bool {
	if u.Role == "admin" {
		return true
	}
	return action == "read"
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
    rows, err := s.repo.database.Query(
      "SELECT id, username, email, role, created_at FROM users",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}
func (s *UserService) DeleteUserByEmail(email string) error {
	_, err := s.repo.database.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}



