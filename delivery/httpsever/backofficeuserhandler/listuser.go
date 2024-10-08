package backofficeuserhandler

import (
	"gameApp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) listUser(c echo.Context) error {
	list, err := h.backofficeuserSvc.ListAllUser()
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": list,
	})
}
