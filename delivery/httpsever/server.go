package httpsever

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpsever/userhandler"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
	"gameApp/vlidator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      config.Config
	userhandler userhandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:      config,
		userhandler: userhandler.New(config.Auth, authSvc, userSvc, userValidator),
	}
}

func (s Server) Serve() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", s.healthCheck)
	s.userhandler.SetUserRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", s.config.HTTPServer.Port)))
}
