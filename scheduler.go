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

func (s *Scheduler) ShortTermSchedule() {
	top := s.ReadyQueue.Dequeue()
	if top != nil {
		top.State = StateRunning
		log.Printf("[Scheduler] Short-term scheduler load task %s to CPU. Task state changed to %s", top.ShortDescription(), top.State)
	}
}

func (s *Scheduler) LongTermSchedule() {
	for len(s.NewQueue.Items()) > 0 {
		t := s.NewQueue.Dequeue()
		s.ReadyQueue.Enqueue(t)
		t.State = StateReady
		log.Printf("[Scheduler] Long-term scheduler moved task %s from %s queue to %s queue. Task state changed to %s\n", t.ShortDescription(), s.NewQueue.Name(), s.ReadyQueue.Name(), t.State)
	}
}
