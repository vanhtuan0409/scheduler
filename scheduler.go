package main

import (
	"log"
	"sync"
)

type TaskQueue struct {
	Name  string
	queue []*Task
	sync.Mutex
}

func NewQueue(name string) *TaskQueue {
	q := new(TaskQueue)
	q.Name = name
	q.queue = []*Task{}
	return q
}

func (q *TaskQueue) Enqueue(t *Task) {
	q.Lock()
	defer q.Unlock()
	log.Printf("[Info] Add task name %s to %s queue", t.Name, q.Name)
	q.queue = append(q.queue, t)
}

type Scheduler struct {
	NewQueue    *TaskQueue
	ReadyQueue  *TaskQueue
	DeviceQueue *TaskQueue //another name is I/O Queue
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.NewQueue = NewQueue("New")
	s.ReadyQueue = NewQueue("Ready")
	s.DeviceQueue = NewQueue("Device")
	return s
}
