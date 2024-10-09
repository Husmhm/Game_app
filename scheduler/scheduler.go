package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
}

func New() Scheduler {
	return Scheduler{}
}

func (sc Scheduler) Start(done <-chan bool) {
	fmt.Println("Scheduler start")
	i := 0
	for {
		select {
		case <-done:
			fmt.Println("exiting ...")
			return
		default:
			now := time.Now()
			fmt.Println("Scheduler now:", now, i)
			time.Sleep(1 * time.Second)
			i++
		}

	}
}
