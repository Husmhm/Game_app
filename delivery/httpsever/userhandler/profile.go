package userhandler

import (
	"gameApp/param"
	"gameApp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userProfile(c echo.Context) error {

	authToken := c.Request().Header.Get("Authorization")
	claims, err := h.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	resp, err := h.userSvc.Profile(param.ProfileRequest{UserID: claims.UserID})

	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
