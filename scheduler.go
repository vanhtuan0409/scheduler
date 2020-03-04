package scheduler

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
