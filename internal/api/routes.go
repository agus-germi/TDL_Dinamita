package api

import (
	//"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *API) SetRoutes(e *echo.Echo) {
	// Public routes (without JWT token)
	public := e.Group("/api/v1")
	public.POST("/auth/signup", a.RegisterUser)
	public.POST("/auth/login", a.LoginUser)

	// Main group of routes for API v1 (with JWT token)
	api := e.Group("/api/v1", a.JWTMiddleware)

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
	tables.GET("", a.GetTables)

	// Group routes for /time_slots under /api/v1
	timeSlots := api.Group("/time_slots")
	timeSlots.GET("", a.GetTimeSlots)
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
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},                    // Origen que permitimos que se conecte a nuestra API.  TODO: Usar variable de entorno para cambiar facilmente el origen entre desarrollo y produccion.
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // A medida que agreguemos mas metodos hay que permitirlos aqui.
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccept}, //para que se puedan enviar headers de cualquier tipo
	}))
}
