package models

import (
	"encoding/json"
	"time"
)

type Transaction struct {
    ID         int       `json:"id"`
    FromUserID int       `json:"from_user_id"`
    ToUserID   int       `json:"to_user_id"`
    Amount     float64   `json:"amount"`
	Type       string    `json:"type"` 
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
}

const (
    TxPending = "pending"
    TxSuccess = "success"
    TxFailed  = "failed"
)

func (t *Transaction) UpdateStatus(s string) { t.Status = s }

func (t *Transaction) ToJSON() ([]byte, error)  { return json.Marshal(t) }
func (t *Transaction) FromJSON(b []byte) error  { return json.Unmarshal(b, t) }
