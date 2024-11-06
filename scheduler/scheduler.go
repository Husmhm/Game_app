package scheduler

import (
	"fmt"
	"gameApp/param"
	"gameApp/service/matchingservice"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Config struct {
	MatchWaitedUsersIntervalInSeconds int `koanf:"match_waited_users_interval_in_seconds"`
}

type Scheduler struct {
	config   Config
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(config Config, matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		config:   config,
		sch:      gocron.NewScheduler(time.UTC),
		matchSvc: matchSvc,
	}
}

func (sc Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	sc.sch.Every(sc.config.MatchWaitedUsersIntervalInSeconds).Second().Do(sc.MatchWaitedUsers)

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
