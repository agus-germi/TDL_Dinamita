package api

import "github.com/labstack/echo/v4"

func (a *API) SetRoutes(e *echo.Echo) {
	users := e.Group("/users")

	users.POST("/register", a.RegisterUser)
}
