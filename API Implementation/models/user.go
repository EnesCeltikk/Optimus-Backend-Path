package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
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

func (u *User) Validate() error {
    if len(u.Username) < 3 {
        return fmt.Errorf("username must be >=3 chars")
    }
    if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`).MatchString(u.Email) {
        return fmt.Errorf("invalid email")
    }
    return nil
}

func (u *User) ToJSON() ([]byte, error)  { return json.Marshal(u) }
func (u *User) FromJSON(b []byte) error  { return json.Unmarshal(b, u) }
