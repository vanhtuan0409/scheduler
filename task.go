package main

type TState int

const (
	StateNew TState = iota
	StateReady
	StateRunning
	StateWaiting
	StateTerminated
)

func (s TState) String() string {
	switch s {
	case StateNew:
		return "New"
	case StateReady:
		return "Ready"
	case StateRunning:
		return "Running"
	case StateWaiting:
		return "Waiting"
	case StateTerminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

type SchedulingInformation struct {
	Priority int
}

type InstructionType int

const (
	CPUBounded InstructionType = iota
	IOBounded
)

func (i InstructionType) String() string {
	switch i {
	case CPUBounded:
		return "CPU Bounded"
	case IOBounded:
		return "I/O Bounded"
	default:
		return "Unknown"
	}
}

type Instruction struct {
	Type     InstructionType
	Duration int
}

type Task struct {
	PID         string
	State       TState
	SInfo       *SchedulingInformation
	Code        []Instruction
	ProgCounter int
}

type TaskTable map[string]*Task
