package scheduler

import (
	"fmt"
	"gameApp/param"
	"gameApp/service/matchingservice"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		sch:      gocron.NewScheduler(time.UTC),
		matchSvc: matchSvc,
	}
}

func (sc Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	sc.sch.Every(5).Second().Do(sc.MatchWaitedUsers)

	sc.sch.StartAsync()
	fmt.Println("Scheduler start")

	<-done
	fmt.Println("stop scheduler ...")
	sc.sch.Stop()

}

func (sc Scheduler) MatchWaitedUsers() {
	fmt.Println("Scheduler MatchWaitedUsers")
	sc.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
