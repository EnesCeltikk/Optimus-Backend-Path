package services

import "CoreImplementation/models"

type TransactionQueue struct {
    C chan *models.Transaction
}

func NewQueue(len int) *TransactionQueue {
    return &TransactionQueue{C: make(chan *models.Transaction, len)}
}

func (q *TransactionQueue) Enqueue(tx *models.Transaction) { q.C <- tx }
func (q *TransactionQueue) Dequeue() *models.Transaction  { return <-q.C }
