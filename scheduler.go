package scheduler

import "log"

type Scheduler struct {
	NewQueue    *FifoQueue
	ReadyQueue  *FifoQueue
	DeviceQueue *FifoQueue //another name is I/O Queue
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.NewQueue = NewFifoQueue("New")
	s.ReadyQueue = NewFifoQueue("Ready")
	s.DeviceQueue = NewFifoQueue("Device")
	return s
}

func (s *Scheduler) ShortTermSelect() *Task {
	top := s.ReadyQueue.Dequeue()
	if top != nil {
		log.Printf("[Scheduler] Short-term scheduler select task %s for running", top.ShortDescription())
	}
	return top
}

func (s *Scheduler) LongTermSelect() *Task {
	top := s.NewQueue.Dequeue()
	if top != nil {
		log.Printf("[Scheduler] Long-term scheduler select task %s for loading into memory", top.ShortDescription())
	}
	return top
}
