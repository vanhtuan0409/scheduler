package scheduler

import "errors"

var (
	ErrMaxPIDReach = errors.New("Max PID reach")
	ErrCPUOccupied = errors.New("CPU is occupied. Cannot load other")
)
