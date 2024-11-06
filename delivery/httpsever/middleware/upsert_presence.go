package middleware

import (
	"gameApp/param"
	"gameApp/pkg/claim"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/timestamp"
	"gameApp/service/presenceservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

// closure
func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			_, err = service.Upsert(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})
			}

			return next(c)
		}
	}
}
