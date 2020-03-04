package scheduler

import (
	"fmt"
	"log"
	"time"
)

var (
	CPUTickInterval = time.Second
)

type CPU struct {
	Instructions []InstructionType
	Timer        *time.Ticker
	progCounter  int
}

func NewCPU() *CPU {
	cpu := new(CPU)
	cpu.Timer = time.NewTicker(CPUTickInterval)
	return cpu
}

func (c *CPU) Report() string {
	instruction := c.fetchNextInstruction()
	if instruction == nil {
		return "CPU is idle"
	} else if *instruction == CPUBounded {
		return fmt.Sprintf("CPU is running. Program Counter: %d", c.progCounter)
	} else if *instruction == IOBounded {
		return fmt.Sprintf("CPU is idle, waiting for I/O. Program Counter: %d", c.progCounter)
	} else {
		return fmt.Sprintf("Error encountered. Exit instruction still occupied CPU. Should already triggered context switch")
	}
}

func (c *CPU) Load(t *Task) error {
	c.Instructions = []InstructionType{}
	for _, i := range t.Code {
		if i.Type == Exit {
			c.Instructions = append(c.Instructions, Exit)
			continue
		}
		for j := 0; j < i.Duration; j++ {
			c.Instructions = append(c.Instructions, i.Type)
		}
	}
	c.progCounter = t.ProgCounter
	return nil
}

func (c *CPU) Work() bool {
	c.progCounter += 1
	instruction := c.fetchNextInstruction()
	if instruction != nil && *instruction == Exit {
		log.Printf("[CPU] Encountered Exit instruction. Return control back to kernel")
		return true
	}
	return false
}

func (c *CPU) fetchNextInstruction() *InstructionType {
	if len(c.Instructions) == 0 || c.progCounter >= len(c.Instructions) {
		return nil
	}
	return &c.Instructions[c.progCounter]
}
