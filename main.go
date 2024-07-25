package main

import (
	"gameApp/config"
	"gameApp/delivery/httpsever"
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
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
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
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
	authService, userService, userValidator := setupServices(cfg)
	// TODO add command for migrations
	//mgr := migrator.New(cfg.Mysql)
	//mgr.Up()

	server := httpsever.New(cfg, authService, userService, userValidator)
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	uV := uservalidator.New(MysqlRepo)

	userSvc := userservice.NewService(authSvc, MysqlRepo)
	return authSvc, userSvc, uV

}
