package middleware

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/claim"
	"gameApp/pkg/errmsg"
	"gameApp/service/authorizationservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

// closure
func AccessCheck(service authorizationservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)
			if err != nil {
				fmt.Println("access control error", isAllowed, err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})
			}

			if !isAllowed {
				fmt.Println("access control !isAllowed", isAllowed, err)
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorMsgUserNotAllowed,
				})
			}
			fmt.Println("access control ok", isAllowed)
			return next(c)
		}
	}
}
