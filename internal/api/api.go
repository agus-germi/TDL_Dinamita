package api

import (
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
	return e.Start(address)
}

func (a *API) Setup(e *echo.Echo) {
	a.SetMiddlewares(e)
	a.SetStaticFiles(e)
	a.SetRoutes(e)
}
