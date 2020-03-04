package scheduler

import "log"

type Scheduler struct {
	NewQueue    Queue
	ReadyQueue  Queue
	DeviceQueue Queue //another name is I/O Queue
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.NewQueue = NewFifoQueue("New")
	s.ReadyQueue = NewFifoQueue("Ready")
	s.DeviceQueue = NewFifoQueue("Device")
	return s
}

func (s *Scheduler) LongTermSchedule() {
	for _, t := range s.NewQueue.Items() {
		s.NewQueue.Dequeue(t)
		s.ReadyQueue.Enqueue(t)
		t.State = StateReady
		log.Printf("[Scheduler] Long-term scheduler moved task pid %d (%s) from %s queue to %s queue. Task state changed to %s\n", t.PID, t.Name, s.NewQueue.Name(), s.ReadyQueue.Name(), t.State)
	}
}
