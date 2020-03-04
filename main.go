package main

import (
	"log"
)

func main() {
	kern := new(Kernel)
	if err := kern.Initialize(); err != nil {
		panic(err)
	}

	t := &Task{
		Name: "Chrome",
		SInfo: &SchedulingInformation{
			Priority: 1,
		},
		Code: []Instruction{
			Instruction{
				Type:     CPUBounded,
				Duration: 5,
			},
		},
	}
	kern.NewTask(t)

	for {
		select {
		case <-kern.Timer.C:
			log.Println("Do scheduling")
		}
	}
}
