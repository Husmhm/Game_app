package httpsever

import (
	"fmt"
	"gameApp/config"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  config.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service) Server {
	return Server{
		config:  config,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) Serve() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health-check", s.healthCheck)

	userGroup := e.Group("/users")

	userGroup.POST("/register", s.userRegister)
	userGroup.POST("/login", s.userLogin)
	userGroup.GET("/profile", s.userProfile)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", s.config.HTTPServer.Port)))
}
