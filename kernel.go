package main

import (
	"log"
	"time"
)

var (
	ClockTickInterval = time.Second
)

type Kernel struct {
	PTable TaskTable
	Timer  *time.Ticker
}

func (k *Kernel) Initialize() error {
	k.PTable = map[int]*Task{}
	initProcess := &Task{Name: "InitV"}
	if err := k.NewTask(initProcess); err != nil {
		return err
	}
	k.Timer = time.NewTicker(ClockTickInterval)
	return nil
}

func (k *Kernel) NewTask(t *Task) error {
	newPID := k.PTable.findSmallestAvailablePID()
	if newPID == -1 {
		return ErrMaxPIDReach
	}
	t.PID = newPID
	t.State = StateNew
	k.PTable[newPID] = t
	log.Printf("[Info] A new process created. Name: %s. PID: %d\n", t.Name, newPID)
	return nil
}
