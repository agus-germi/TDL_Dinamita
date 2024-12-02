package api

import (
	//"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO: Cambiar las URIs de los endpoints para que sean mas RESTful.
func (a *API) SetRoutes(e *echo.Echo) {
	// Public routes (without JWT token)
	public := e.Group("/api/v1")
	public.POST("/users", a.RegisterUser)
	public.POST("/users/login", a.LoginUser)

	// Main group of routes for API v1 (with JWT token)
	api := e.Group("/api/v1", JWTMiddleware)

	// Group routes for /users under /api/v1
	users := api.Group("/users")
	users.DELETE("/:id", a.DeleteUser)
	users.PATCH("/:id", a.UpdateUserRole)
	users.GET("/:id/reservations", a.GetReservationsOfUser)

	// Group routes for /reservations under /api/v1
	reservations := api.Group("/reservations")
	reservations.POST("", a.CreateReservation)
	reservations.DELETE("/:id", a.DeleteReservation)

	// Group routes for /tables under /api/v1
	tables := api.Group("/tables")
	tables.POST("", a.CreateTable)
	tables.DELETE("/:id", a.DeleteTable)
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
		AllowOrigins:     []string{"http://localhost:8080"},                    // Origen que permitimos que se conecte a nuestra API.  TODO: Usar variable de entorno para cambiar facilmente el origen entre desarrollo y produccion.
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // A medida que agreguemos mas metodos hay que permitirlos aqui.
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccept}, //para que se puedan enviar headers de cualquier tipo
	}))
}
