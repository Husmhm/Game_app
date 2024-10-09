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
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator,
	authorizationSvc authorizationservice.Service, backofficeuserSvc backofficeuserservice.Service,
	matchSvc matchingservice.Service, matchvalidator matchingvalidator.Validator) Server {
	return Server{
		config:                config,
		userhandler:           userhandler.New(config.Auth, authSvc, userSvc, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc, authorizationSvc, backofficeuserSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSvc, matchSvc, matchvalidator),
	}
}

func (s Server) Serve() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", s.healthCheck)
	s.userhandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)
	s.matchingHandler.SetRoutes(e)

	// Start server
	address := fmt.Sprintf("localhost:%d", s.config.HTTPServer.Port)
	fmt.Printf("star echo server on %s\n", address)
	e.Logger.Fatal(e.Start(address))
}
