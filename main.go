package main

import (
	"fmt"
	"gameApp/adapter/redis"
	"gameApp/config"
	"gameApp/delivery/httpsever"
	"gameApp/repository/migrator"
	"gameApp/repository/mysql"
	"gameApp/repository/mysql/mysqlaccesscontrol"
	"gameApp/repository/mysql/mysqluser"
	"gameApp/repository/redis/redismatching"
	"gameApp/repository/redis/redispresence"
	"gameApp/scheduler"
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/matchingservice"
	"gameApp/service/presenceservice"
	"gameApp/service/userservice"
	"gameApp/vlidator/matchingvalidator"
	"gameApp/vlidator/uservalidator"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	// TODO - read config path from command line

	cfg := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg)

	// TODO -add struct and add these returned items struct field.
	authService, userService, userValidator, authorizationSvc, backOfficeUserSvc, matchingSvc, matchingV, presenceSvc := setupServices(cfg)

	// TODO add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	server := httpsever.New(cfg, authService, userService, userValidator, authorizationSvc, backOfficeUserSvc, matchingSvc, matchingV, presenceSvc)
	go func() {
		server.Serve()
	}()

	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(cfg.Schduller, matchingSvc)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GraceFullShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctx); err != nil {
		fmt.Println("http server shutdown error", err)
	}
	fmt.Println("\n recieved interrupt signal, shutting down gracefully...")
	done <- true
	time.Sleep(cfg.Application.GraceFullShutdownTimeout)

	// TODO- does order of ctx.done &wg.wait matter?
	ctx.Done()

	wg.Wait()
	time.Sleep(2 * time.Second)
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator,
	authorizationservice.Service, backofficeuserservice.Service,
	matchingservice.Service, matchingvalidator.Validator, presenceservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(MysqlRepo)
	uV := uservalidator.New(userMysql)
	userSvc := userservice.NewService(authSvc, userMysql)

	backOfficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	matchingV := matchingvalidator.New()

	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	//TODO - panic - replace presenceSvc with presence grpc client
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceSvc)

	return authSvc, userSvc, uV, authorizationSvc, backOfficeUserSvc, matchingSvc, matchingV, presenceSvc

}
