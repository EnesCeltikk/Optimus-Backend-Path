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

func (s *UserService) Login(username, password string) (string, error) {
	u, err := s.repo.FindByEmail(username)
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	// Dummy token
	return "dummy-token", nil
}

func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	// Dummy refresh
	if refreshToken == "dummy-refresh" {
		return "new-dummy-token", nil
	}
	return "", errors.New("invalid refresh token")
}

func (s *UserService) UpdateUser(user *models.User) error {
	_, err := s.repo.database.Exec("UPDATE users SET username=$1, email=$2, role=$3 WHERE id=$4", user.Username, user.Email, user.Role, user.ID)
	return err
}

func (s *UserService) DeleteUser(id int) error {
	_, err := s.repo.database.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	row := s.repo.database.QueryRow("SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id=$1", id)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

type Claims struct {
	UserID int
	Role   string
}

// Dummy implementation for token validation
func ValidateToken(token string) (*Claims, error) {
	if token == "dummy-token" {
		return &Claims{UserID: 1, Role: "admin"}, nil
	}
	return nil, errors.New("invalid token")
}



