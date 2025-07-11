package models

type UserRepository interface {
    Insert(u *User) error
    FindByEmail(email string) (*User, error)
}

type TransactionRepository interface {
    Create(tx *Transaction) error
    UpdateStatus(id int, status string) error
}

type BalanceRepository interface {
    Get(userID int) (*Balance, error)
    Update(b *Balance) error
    History(userID int) ([]Balance, error)
}

type UserService interface {
    Register(username, email, password string) (*User, error)
    Authenticate(email, password string) (*User, error)
    Authorize(user *User, action string) bool
}

type TransactionService interface {
    Transfer(fromID, toID int, amt float64) (*Transaction, error)
    Rollback(txID int) error
}

type BalanceService interface {
    Update(userID int, amt float64) error
    Get(userID int) (*Balance, error)
    History(userID int) ([]Balance, error)
}
