package backofficeuserhandler

import (
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizationservice.Service
	backofficeuserSvc backofficeuserservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service, authorizationSvc authorizationservice.Service, backofficeuserSvc backofficeuserservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		authorizationSvc:  authorizationSvc,
		backofficeuserSvc: backofficeuserSvc,
	}
}
