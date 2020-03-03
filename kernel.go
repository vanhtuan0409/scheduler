package main

import "log"

type Kernel struct {
	PTable TaskTable
}

func (k *Kernel) Initialize() error {
	k.PTable = map[int]*Task{}
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
	log.Printf("[Info] A new process created. PID: %d\n", newPID)
	return nil
}
