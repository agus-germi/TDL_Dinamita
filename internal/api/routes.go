package api

import (
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *API) SetRoutes(e *echo.Echo) {
	//api := e.Group("/api")

	// Group routes for /users
	users := e.Group("/users")
	users.POST("/register", a.RegisterUser)

	// Set static files
	//a.setStaticFiles(e)
}

// Aca hay que definir bien como estructuramos el directorio "static".
// Podemos crear una carpeta por recurso y hacer la jerarquia de carpetas de forma
// que con setear los StaticFiles como 'e.Static("/users", "static/users")'
// la busqueda se haga automatica.
func (a *API) SetStaticFiles(e *echo.Echo) {
	//publicDir := "frontend" // Just use the folder name since it will be relative to the working directory

	// Serve the index.html at the root
	e.File("/", "index.html")

	// Serve all static assets (CSS, JS, etc.) from "/static" prefix
	e.Static("/", "frontend") // This serves the frontend directory at the /static path
}


func (a *API) SetMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
