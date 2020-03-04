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
		currInstruction := c.findCurrentInstruction()
		if currInstruction.Type == CPUBounded {
			return fmt.Sprintf("CPU is running. Task %s. Program Counter: %d/%d", c.RunningTask.ShortDescription(), c.progCounter, c.RunningTask.TotalDuration())
		} else {
			return fmt.Sprintf("CPU is idle, waiting for I/O. Task %s. Program Counter: %d/%d", c.RunningTask.ShortDescription(), c.progCounter, c.RunningTask.TotalDuration())
		}
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

func (c *CPU) findCurrentInstruction() *Instruction {
	counter := c.progCounter
	for _, i := range c.RunningTask.Code {
		if counter > i.Duration {
			counter -= i.Duration
		} else {
			return &i
		}
	}
	return &c.RunningTask.Code[len(c.RunningTask.Code)-1]
}
