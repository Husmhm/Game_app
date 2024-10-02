package backofficeuserhandler

import (
	"gameApp/delivery/httpsever/middleware"
	"gameApp/entity"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/", h.listUser, middleware.Auth(h.authSvc, h.authConfig),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))
}
