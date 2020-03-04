package scheduler

import "fmt"

const (
	MaxPID = 4194303 // equal to config of /proc/sys/kernel/pid_max
)

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
	Exit
)

func (i InstructionType) String() string {
	switch i {
	case CPUBounded:
		return "CPU Bounded"
	case IOBounded:
		return "I/O Bounded"
	case Exit:
		return "Exit"
	default:
		return "Unknown"
	}
}

type Instruction struct {
	Type     InstructionType
	Duration int
}

type Task struct {
	PID         int
	Name        string //In real life, this should be the bootstrap command of this task
	State       TState
	SInfo       *SchedulingInformation
	Code        []Instruction
	ProgCounter int
}

func (t *Task) ShortDescription() string {
	return fmt.Sprintf("pid %d (%s)", t.PID, t.Name)
}

func (t *Task) TotalDuration() int {
	res := 0
	for _, i := range t.Code {
		res += i.Duration
	}
	return res
}

func (t *Task) IsFinished() bool {
	return t.ProgCounter >= t.TotalDuration()
}

type TaskTable map[int]*Task

func (tab TaskTable) findSmallestAvailablePID() int {
	for i := 0; i < MaxPID; i++ {
		if tab[i] == nil {
			return i
		}
	}
	return -1
}
