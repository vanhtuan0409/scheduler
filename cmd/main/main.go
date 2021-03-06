package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vanhtuan0409/scheduler"
)

func main() {
	// Initialize Kernel
	kern := new(scheduler.Kernel)
	if err := kern.Initialize(); err != nil {
		panic(err)
	}
	// Set option for kernel. In real life, this is usually done via `sysctl`
	kern.Options.DisableLongTermScheduler = false

	// Add simple task
	t := &scheduler.Task{
		Name: "Chrome",
		SInfo: &scheduler.SchedulingInformation{
			Priority: 1,
		},
		Code: []scheduler.Instruction{
			scheduler.Instruction{
				Type:     scheduler.CPUBounded,
				Duration: 2,
			},
			scheduler.Instruction{
				Type:     scheduler.IOBounded,
				Duration: 2,
			},
			scheduler.Instruction{
				Type: scheduler.Exit,
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
		log.Println("[Info] Receive halting signal")
		kern.Halt()
	}()

	// Loop for scheduling activity
	// In real life, scheduler is wake up by an interrupt
	// (eg: disk done, mouse click, timer tick)
	// During the timer, CPU will be doing its job without interfere
schedulingLoop:
	for {
		select {

		case <-kern.Core1.Timer.C:
			if kern.Core1.Work() {
				kern.CleanupRunningTask()
			}

		case <-kern.ShortTermScheduleTimer.C:
			kern.DoShortTermScheduling()

		case <-kern.LongTermScheduleTimer.C:
			kern.DoLongTermScheduling()

		case <-kern.Exited():
			break schedulingLoop
		}
	}

	log.Println("[Info] System shutdown")
}
