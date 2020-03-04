package scheduler

import (
	"log"
	"time"
)

var (
	LongTermInterval  = 5 * time.Second
	ShortTermInterval = 2 * time.Second
)

type KernelOptions struct {
	DisableLongTermScheduler bool
	PreEmptive               bool
}

type Kernel struct {
	PTable                 TaskTable
	ShortTermScheduleTimer *time.Ticker
	LongTermScheduleTimer  *time.Ticker
	Scheduler              *Scheduler
	Options                KernelOptions
	exitChan               chan struct{}
}

func (k *Kernel) Initialize() error {
	k.PTable = map[int]*Task{}
	initProcess := &Task{Name: "InitV"}
	if err := k.NewTask(initProcess); err != nil {
		return err
	}
	k.Scheduler = NewScheduler()
	k.ShortTermScheduleTimer = time.NewTicker(ShortTermInterval)
	k.LongTermScheduleTimer = time.NewTicker(LongTermInterval)
	k.exitChan = make(chan struct{})
	return nil
}

func (k *Kernel) Halt() {
	// do halting logic here
	log.Println("[Kern] Prepare for halting")
	k.exitChan <- struct{}{}
}

func (k *Kernel) Exited() <-chan struct{} {
	return k.exitChan
}

func (k *Kernel) NewTask(t *Task) error {
	newPID := k.PTable.findSmallestAvailablePID()
	if newPID == -1 {
		return ErrMaxPIDReach
	}
	t.PID = newPID
	t.State = StateNew
	k.PTable[newPID] = t
	log.Printf("[Kern] A new process created. Name: %s. PID: %d\n", t.Name, newPID)
	if newPID == 0 {
		return nil
	}

	// Perform queueing task to scheduler
	if k.Options.DisableLongTermScheduler {
		k.Scheduler.ReadyQueue.Enqueue(t)
	} else {
		k.Scheduler.NewQueue.Enqueue(t)
	}
	return nil
}

func (k *Kernel) CleanupTask(t *Task) error {
	t.State = StateTerminated
	log.Printf("[Kern] Process %s finished. Cleaning its state\n", t.ShortDescription())
	delete(k.PTable, t.PID)
	return nil
}
