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

	reservations := e.Group("/reservations")
	reservations.POST("/register", a.RegisterReservation)
	reservations.DELETE("/remove", a.RemoveReservation)

	tables := e.Group("/tables")
	tables.POST("/register", a.AddTable)
	tables.DELETE("/remove", a.RemoveTable)

	// Handle GET requests to the root path
	//e.GET("/*", func(c echo.Context) error {
	//	return c.File("static/index.html")
	//})

	// Set static files
	//a.setStaticFiles(e)
}

// Aca hay que definir bien como estructuramos el directorio "static".
// Podemos crear una carpeta por recurso y hacer la jerarquia de carpetas de forma
// que con setear los StaticFiles como 'e.Static("/users", "static/users")'
// la busqueda se haga automatica.
func (a *API) SetStaticFiles(e *echo.Echo) {
	// Determina la ruta absoluta al directorio "public"
	publicDir := filepath.Join("frontend", "public")

	// Serve all static assets (CSS, JS, etc.) from "/static" prefix
	e.Static("/", publicDir)

	// Serve the index.html at the root
	e.File("/", filepath.Join(publicDir, "index.html"))
}

func (a *API) SetMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
