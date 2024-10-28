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
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/matchingservice"
	"gameApp/service/userservice"
	"gameApp/vlidator/matchingvalidator"
	"gameApp/vlidator/uservalidator"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"os"
	"os/signal"
)

func main() {
	// TODO - read config path from command line

	cfg := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg)

	// TODO -add struct and add these returned items struct field.
	authService, userService, userValidator, authorizationSvc, backOfficeUserSvc, mathingSvc, matchingV := setupServices(cfg)
	// TODO add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	var httpServer *echo.Echo
	go func() {
		server := httpsever.New(cfg, authService, userService, userValidator, authorizationSvc, backOfficeUserSvc, mathingSvc, matchingV)
		httpServer = server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GraceFullShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Println("http server shutdown error", err)
	}
	fmt.Println("\n recieved interrupt signal, shutting down gracefully...")

	<-ctx.Done()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator,
	authorizationservice.Service, backofficeuserservice.Service,
	matchingservice.Service, matchingvalidator.Validator) {
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
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo)

	return authSvc, userSvc, uV, authorizationSvc, backOfficeUserSvc, matchingSvc, matchingV

}
