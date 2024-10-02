package claim

import (
	"gameApp/config"
	"gameApp/service/authservice"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	return c.Get(config.AuthMiddleWareConrexKey).(*authservice.Claims)
}
