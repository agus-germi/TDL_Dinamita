package api

import (
	"log"
	"net/http"

	"github.com/agus-germi/TDL_Dinamita/jwtutils"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := jwtutils.GetClaimsFromCookie(c)
		if err != nil {
			log.Println("Error while getting the JWT claims from the cookie:", err)
			return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
		}

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		return next(c)
	}
}
