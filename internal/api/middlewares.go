package api

import (
	"net/http"

	"github.com/agus-germi/TDL_Dinamita/jwtutils"
	"github.com/labstack/echo/v4"
)

func (a *API) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.log.Debugf("[Middleware] called")

		claims, err := jwtutils.GetClaimsFromToken(c)
		if err != nil {
			a.log.Errorf("Error while getting the JWT claims:", err)
			return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
		}

		a.log.Debugf("[Middleware] Claims:", claims)
		a.log.Debugf("[Middleware] usr_id:", claims["user_id"])
		a.log.Debugf("[Middleware] email:", claims["email"])
		a.log.Debugf("[Middleware] role_id:", claims["role_id"])

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role_id"])

		a.log.Debugf("[Middleware] GET usr_id:", c.Get("user_id"))
		a.log.Debugf("[Middleware] GET email:", c.Get("email"))
		a.log.Debugf("[Middleware] GET role_id:", c.Get("role"))

		a.log.Debugf("[Middleware] finished successfully")
		return next(c)
	}
}

func (a *API) validateAdminRole(c echo.Context) (int64, error) {
	clientRoleID, ok := c.Get("role").(float64)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "[Middleware] Invalid client role ID in context")
	}
	clientRoleIDInt := int64(clientRoleID)
	if clientRoleIDInt != adminRoleID {
		return 0, echo.NewHTTPError(http.StatusForbidden, "[Middleware] Permission denied: admin role required")
	}
	return clientRoleIDInt, nil
}
