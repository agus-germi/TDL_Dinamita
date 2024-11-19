package api

import (
	//"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *API) SetRoutes(e *echo.Echo) {
	//api := e.Group("/api")

	// Group routes for /users
	users := e.Group("/users")
	users.POST("/register", a.RegisterUser)
	users.DELETE("/remove", a.RemoveUser)
	users.POST("/roles", a.AddUserRole)
	users.POST("/login", a.LoginUser)

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
	// Sirve todos los archivos estáticos desde el prefijo /static
	e.Static("/static", "frontend")

	// Sirve el archivo index.html en la raíz
	e.File("/", "frontend/index.html")
}

func (a *API) SetMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},          // Origen que permitimos que se conecte a nuestra API.
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE}, // A medida que agreguemos mas metodos hay que permitirlos aqui.
		AllowCredentials: true,
		// AllowHeaders: []string{echo.HeaderContentType},  // Analizar si queremos restringir los tipos de headers.
	}))
}
