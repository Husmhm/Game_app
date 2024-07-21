package main

import (
	"gameApp/config"
	"gameApp/delivery/httpsever"
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
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
	authService, userService := setupServices(cfg)

	server := httpsever.New(cfg, authService, userService)
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.NewService(authSvc, MysqlRepo)
	return authSvc, userSvc

}
