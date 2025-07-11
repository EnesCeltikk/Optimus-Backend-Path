package models

import (
	"encoding/json"
	"sync"
	"time"
)

type Balance struct {
    UserID        int       `json:"user_id"`
    Amount        float64   `json:"amount"`
    LastUpdatedAt time.Time `json:"last_updated_at"`
}

var (
    balanceMutex   = make(map[int]*sync.RWMutex)
    balanceHistory = make(map[int][]Balance)
    histMu         sync.Mutex
)

func getMu(userID int) *sync.RWMutex {
    histMu.Lock()
    defer histMu.Unlock()
    if _, ok := balanceMutex[userID]; !ok {
        balanceMutex[userID] = &sync.RWMutex{}
    }
    return balanceMutex[userID]
}

func (b *Balance) Update(amount float64) {
    mu := getMu(b.UserID)
    mu.Lock()
    defer mu.Unlock()
    b.Amount += amount
    b.LastUpdatedAt = time.Now()

    histMu.Lock()
    balanceHistory[b.UserID] = append(balanceHistory[b.UserID], *b)
    histMu.Unlock()
}

func GetHistory(userID int) []Balance {
    histMu.Lock()
    defer histMu.Unlock()
    return balanceHistory[userID]
}

func (b *Balance) ToJSON() ([]byte, error)  { return json.Marshal(b) }
func (b *Balance) FromJSON(d []byte) error  { return json.Unmarshal(d, b) }
