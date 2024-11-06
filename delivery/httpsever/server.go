package httpsever

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpsever/backofficeuserhandler"
	"gameApp/delivery/httpsever/matchinghandler"
	"gameApp/delivery/httpsever/userhandler"
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/matchingservice"
	"gameApp/service/presenceservice"
	"gameApp/service/userservice"
	"gameApp/vlidator/matchingvalidator"
	"gameApp/vlidator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userhandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
	Router                *echo.Echo
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator,
	authorizationSvc authorizationservice.Service, backofficeuserSvc backofficeuserservice.Service,
	matchSvc matchingservice.Service, matchvalidator matchingvalidator.Validator, presenceSvc presenceservice.Service) Server {
	return Server{
		config:                config,
		userhandler:           userhandler.New(config.Auth, authSvc, userSvc, userValidator, presenceSvc),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc, authorizationSvc, backofficeuserSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSvc, matchSvc, matchvalidator, presenceSvc),
		Router:                echo.New(),
	}
}

func (s Server) Serve() {
	// Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// Routes
	s.Router.GET("/health-check", s.healthCheck)
	s.userhandler.SetRoutes(s.Router)
	s.backofficeUserHandler.SetRoutes(s.Router)
	s.matchingHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf("localhost:%d", s.config.HTTPServer.Port)
	fmt.Printf("star echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Printf("router start server error: %v\n", err)
	}

}
