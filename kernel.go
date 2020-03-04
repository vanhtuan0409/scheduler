package scheduler

import (
	"bytes"
	"fmt"
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
	RunningTask            *Task
	Options                KernelOptions
	exitChan               chan struct{}
	Core1                  *CPU // in real life, Kernel should be able to communicate directly to CPU
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

	// Setting up device driver
	k.Core1 = NewCPU()
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

func (k *Kernel) CleanupRunningTask() error {
	t := k.RunningTask
	t.State = StateTerminated
	log.Printf("[Kern] Process %s finished. Cleaning its state\n", t.ShortDescription())
	delete(k.PTable, t.PID)

	nextTask := k.Scheduler.ShortTermSchedule()
	if nextTask != nil {
		k.ContextSwitch(nextTask)
	}
	return nil
}

func (k *Kernel) ContextSwitch(t *Task) error {
	buf := bytes.NewBufferString("[Kern] Performing context switch")
	if !k.IsCPUFree() {
		buf.WriteString(fmt.Sprintf(". Swapping out process %s", k.RunningTask.ShortDescription()))
		k.SwapOut(k.RunningTask)
	}
	buf.WriteString(fmt.Sprintf(". Swapping in process %s", t.ShortDescription()))
	k.SwapIn(t)
	log.Printf("%s\n", buf.String())
	return nil
}

func (k *Kernel) SwapIn(t *Task) {
	k.RunningTask = t
	t.State = StateRunning
	k.Core1.Load(t)
}

func (k *Kernel) SwapOut(t *Task) {
	// Save register
	t.ProgCounter = k.Core1.progCounter
	k.RunningTask = nil
	// Put it to approriate queue
	// TODO
}

func (k *Kernel) IsCPUFree() bool {
	return k.RunningTask == nil
}
