package domain

import (
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.PasswordHash == "" {
		return errors.New("password is required")
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
} 