package scheduler

import (
	"fmt"
	"time"
)

var (
	CPUTickInterval = time.Second
)

type CPU struct {
	RunningTask *Task
	Timer       *time.Ticker
	progCounter int
}

func NewCPU() *CPU {
	cpu := new(CPU)
	cpu.Timer = time.NewTicker(CPUTickInterval)
	return cpu
}

func (c *CPU) Report() string {
	if !c.IsFree() {
	}
	return fmt.Sprintf("CPU is free. No occupied task")
}

func (c *CPU) Load(t *Task) error {
	if c.RunningTask != nil {
		return ErrCPUOccupied
	}
	c.RunningTask = t
	c.progCounter = t.ProgCounter
	return nil
}

func (c *CPU) Work() {
	c.progCounter += 1
}

func (c *CPU) IsFree() bool {
	return c.RunningTask == nil
}
