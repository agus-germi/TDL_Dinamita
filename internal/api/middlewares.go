package api

import (
	"net/http"

	"github.com/agus-germi/TDL_Dinamita/jwtutils"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates the JWT token present in the request header
// and extracts its claims. If the token is valid, the user data (such as
// user_id, email, and role_id) is added to the context for use in subsequent
// handlers. If the token is invalid, a 401 Unauthorized error is returned.
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

// ValidateAdminRole is a middleware that validates the user's role.
// It ensures that the user has admin privileges before allowing access
// to the next handler. If the user does not have admin privileges, it
// returns a 403 Forbidden response.
func (a *API) ValidateAdminRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.log.Debugf("[Middleware] called")

		clientRoleID, ok := c.Get("role").(float64)
		if !ok {
			a.log.Errorf("[Middleware] Invalid client role ID in context")
			return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
		}

		clientRoleIDInt := int64(clientRoleID)
		a.log.Debugf("[Middleware] Role ID:", clientRoleIDInt)

		if clientRoleIDInt != adminRoleID {
			a.log.Errorf("[Middleware] Permission denied: admin role required")
			return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: admin role required"})
		}

		a.log.Debugf("[Middleware] Role validated successfully")
		return next(c)
	}
}
