package api

import (
	"errors"
	"net/http"

	"github.com/agus-germi/TDL_Dinamita/internal/api/dtos"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/labstack/echo/v4"
)

func (a *API) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	params := dtos.RegisterUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = a.serv.RegisterUser(ctx, params.Name, params.Password, params.Email)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, errors.New("unexpected error"))
	}

	return c.JSON(http.StatusCreated, nil)
}
