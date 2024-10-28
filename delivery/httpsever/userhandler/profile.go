package userhandler

import (
	"gameApp/param"
	"gameApp/pkg/claim"
	"gameApp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func (h Handler) userProfile(c echo.Context) error {

	claims := claim.GetClaimsFromEchoContext(c)

	ctx := c.Request().Context()
	ctxWithTimeout, cancle := context.WithTimeout(ctx, time.Second)
	defer cancle()

	resp, err := h.userSvc.Profile(ctxWithTimeout, param.ProfileRequest{UserID: claims.UserID})

	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
