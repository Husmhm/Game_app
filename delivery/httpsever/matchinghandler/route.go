package matchinghandler

import (
	"gameApp/delivery/httpsever/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/matching")

	userGroup.POST("/add-to-wating-list", h.addToWatingList, middleware.Auth(h.authSvc, h.authConfig))

}
