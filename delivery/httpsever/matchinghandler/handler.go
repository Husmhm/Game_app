package matchinghandler

import (
	"gameApp/service/authservice"
	"gameApp/service/matchingservice"
	"gameApp/service/presenceservice"
	"gameApp/vlidator/matchingvalidator"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingSvc       matchingservice.Service
	matchingValidator matchingvalidator.Validator
	presenceSvc       presenceservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service, matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		matchingSvc:       matchingSvc,
		matchingValidator: matchingValidator,
		presenceSvc:       presenceSvc,
	}
}
