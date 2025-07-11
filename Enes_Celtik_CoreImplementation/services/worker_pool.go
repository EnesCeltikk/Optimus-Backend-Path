package services

import (
	"log"
	"sync/atomic"
	"time"

	"CoreImplementation/models"
)

type WorkerPool struct {
    Queue   chan *models.Transaction
    Workers int
    totalProcessed uint64
}

func NewWorkerPool(n, qlen int) *WorkerPool {
    return &WorkerPool{
        Queue:   make(chan *models.Transaction, qlen),
        Workers: n,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.Workers; i++ {
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    for tx := range wp.Queue {
        log.Printf("[worker %d] processing tx %d", id, tx.ID)

        time.Sleep(100 * time.Millisecond)
        tx.UpdateStatus(models.TxSuccess)
        atomic.AddUint64(&wp.totalProcessed, 1)
    }
}

func (wp *WorkerPool) Stats() uint64 {
    return atomic.LoadUint64(&wp.totalProcessed)
}
