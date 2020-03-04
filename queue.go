package scheduler

import (
	"log"
	"sync"
)

type Queue interface {
	Name() string
	Enqueue(*Task)
}

type NamedQueue struct {
	name string
}

func (q NamedQueue) Name() string {
	return q.name
}

type FifoQueue struct {
	NamedQueue
	queue []*Task
	sync.Mutex
}

func NewFifoQueue(name string) *FifoQueue {
	q := new(FifoQueue)
	q.name = name
	q.queue = []*Task{}
	return q
}

func (q *FifoQueue) Enqueue(t *Task) {
	q.Lock()
	defer q.Unlock()
	log.Printf("[Info] Add task name %s to %s queue", t.Name, q.Name())
	q.queue = append(q.queue, t)
}
