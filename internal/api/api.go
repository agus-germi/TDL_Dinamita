package api

import (
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type API struct {
	serv          service.Service
	dataValidator *validator.Validate
	log           logger.Logger
}

func New(serv service.Service, dataValidator *validator.Validate, log logger.Logger) *API {
	log.Debugf("Logger has been injected into API")
	return &API{
		serv:          serv,
		dataValidator: dataValidator,
		log:           log,
	}
}

func (a *API) Start(e *echo.Echo, address string) error {
	a.SetMiddlewares(e)
	a.SetStaticFiles(e)
	a.SetRoutes(e)
	return e.Start(address)
}
