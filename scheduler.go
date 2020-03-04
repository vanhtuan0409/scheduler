package main

import "log"

type TaskQueue struct {
	Name  string
	queue []*Task
}

func NewQueue(name string) *TaskQueue {
	q := new(TaskQueue)
	q.Name = name
	q.queue = []*Task{}
	return q
}

func (q *TaskQueue) Enqueue(t *Task) {
	log.Printf("[Info] Add task name %s to %s queue", t.Name, q.Name)
	q.queue = append(q.queue, t)
}

type Scheduler struct {
	NewQueue   *TaskQueue
	ReadyQueue *TaskQueue
	IOQueue    *TaskQueue
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.NewQueue = NewQueue("New")
	s.ReadyQueue = NewQueue("Ready")
	s.IOQueue = NewQueue("I/O")
	return s
}
