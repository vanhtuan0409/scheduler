package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vanhtuan0409/scheduler"
)

func main() {
	// Initialize CPU
	core1 := scheduler.NewCPU()

	// Initialize Kernel
	kern := new(scheduler.Kernel)
	if err := kern.Initialize(); err != nil {
		panic(err)
	}
	// Set option for kernel. In real life, this is usually done via `sysctl`
	kern.Options.DisableLongTermScheduler = true

	// Add simple task
	t := &scheduler.Task{
		Name: "Chrome",
		SInfo: &scheduler.SchedulingInformation{
			Priority: 1,
		},
		Code: []scheduler.Instruction{
			scheduler.Instruction{
				Type:     scheduler.CPUBounded,
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
		case <-core1.Timer.C:
			log.Printf("[CPU] %s\n", core1.Report())
		case <-kern.ShortTermScheduleTimer.C:
			log.Println("[Info] Short-term scheduler is woke up. Do scheduling")
		case <-kern.LongTermScheduleTimer.C:
			if !kern.Options.DisableLongTermScheduler {
				log.Println("[Info] Long-term scheduler is woke up. Do scheduling")
				kern.Scheduler.LongTermSchedule()
			}
		case <-kern.Exited():
			break schedulingLoop
		}
	}

	log.Println("System shutdown")
}
