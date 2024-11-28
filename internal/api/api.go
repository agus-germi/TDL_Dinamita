package api

import (
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	serv          service.Service
	dataValidator *validator.Validate
}

func New(serv service.Service) *API {
	return &API{
		serv:          serv,
		dataValidator: validator.New(),
	}
}

func (a *API) Start(e *echo.Echo, address string) error {

	a.SetRoutes(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                                          //para que cualquier origen pueda hacer peticiones
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},                   //para que se puedan hacer peticiones de cualquier tipo
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}, //para que se puedan enviar headers de cualquier tipo
	}))

	return e.Start(address)
}
