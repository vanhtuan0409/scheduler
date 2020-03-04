package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// Trapping OS signal to perform halting
	// In real life, halting signal should come from hardware
	// (whenever you click the power button)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)

	go func() {
		<-c
		log.Println("Receive halting signal")
		kern.Halt()
	}()

	// Loop for scheduling activity
	// In real life, scheduler is wake up by an interrupt
	// (eg: disk done, mouse click, timer tick)
	// During the timer, CPU will be doing its job without interfere
schedulingLoop:
	for {
		select {
		case <-kern.Timer.C:
			log.Println("Scheduler is wake up. Do scheduling")
		case <-kern.Exited():
			break schedulingLoop
		}
	}

	log.Println("System shutdown")
}
