package main

func main() {
	kern := new(Kernel)
	if err := kern.Initialize(); err != nil {
		panic(err)
	}

	t := &Task{
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
}
