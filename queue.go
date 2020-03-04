package scheduler

import (
	"log"
	"sync"
)

type Queue interface {
	Name() string
	Enqueue(*Task)
	Dequeue(*Task)
	Items() []*Task
}

type NamedQueue struct {
	name string
}

func (q NamedQueue) Name() string {
	return q.name
}

// Implement FIFO Queue
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
	log.Printf("[Scheduler] Add task %s to %s queue", t.ShortDescription(), q.Name())
	q.queue = append(q.queue, t)
}

func (q *FifoQueue) Dequeue(t *Task) {
	tIndex := -1
	for idx, i := range q.queue {
		if i.PID == t.PID {
			tIndex = idx
		}
	}
	q.queue = append(q.queue[:tIndex], q.queue[tIndex+1:]...)
}

func (q *FifoQueue) Items() []*Task {
	return q.queue
}
