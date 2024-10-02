package main

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpsever"
	"gameApp/repository/migrator"
	"gameApp/repository/mysql"
	"gameApp/repository/mysql/mysqlaccesscontrol"
	"gameApp/repository/mysql/mysqluser"
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/userservice"
	"gameApp/vlidator/uservalidator"
	"time"
)

const (
	Jwt_SignKey                = "jwt_secret"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
)

func main() {
	// TODO - read config path from command line

	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg2)

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey:               Jwt_SignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessTokenSubject:    AccessTokenSubject,
			RefreshTokenSubject:   RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			UserName: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}
	authService, userService, userValidator, authorizationSvc, backOfficeUserSvc := setupServices(cfg)
	// TODO add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	server := httpsever.New(cfg, authService, userService, userValidator, authorizationSvc, backOfficeUserSvc)
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator,
	authorizationservice.Service, backofficeuserservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(MysqlRepo)
	uV := uservalidator.New(userMysql)
	userSvc := userservice.NewService(authSvc, userMysql)

	backOfficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	return authSvc, userSvc, uV, authorizationSvc, backOfficeUserSvc

}
