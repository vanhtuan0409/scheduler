package scheduler

import (
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
	instruction := c.fetchNextInstruction()

	if instruction == nil {
		log.Printf("[CPU] CPU is idle\n")
		return false
	} else if *instruction == CPUBounded {
		log.Printf("[CPU] CPU is running. Program Counter: %d\n", c.progCounter)
		return false
	} else if *instruction == IOBounded {
		log.Printf("[CPU] CPU is idle, waiting for I/O. Program Counter: %d\n", c.progCounter)
		return false
	} else if *instruction == Exit {
		log.Printf("[CPU] Encountered Exit instruction. Return control back to kernel\n")
		return true
	} else {
		log.Printf("[CPU] Encountered Unknown instruction: %s\n", *instruction)
		return true
	}
}

func (c *CPU) fetchNextInstruction() *InstructionType {
	var instruction *InstructionType
	if len(c.Instructions) > 0 && c.progCounter < len(c.Instructions) {
		instruction = &c.Instructions[c.progCounter]
	}
	c.progCounter++
	return instruction
}
