package api

import (
	"log"
	"net/http"

	"github.com/agus-germi/TDL_Dinamita/internal/api/dtos"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/labstack/echo/v4"
)

type responseMessage struct {
	Message string `json:"message"`
}

func (a *API) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	params := dtos.RegisterUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"}) // otra opcion--> Message: err.Error()
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.RegisterUser(ctx, params.Name, params.Password, params.Email)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "User registered successfully"})
}

func (a *API) RegisterReservation(c echo.Context) error {

	ctx := c.Request().Context() // obtengo el contexto del objeto Request que viene con la petición HTTP

	params := dtos.RegisterReservationDTO{} //creo una instancia de RegisterReservationDTO

	//Linkeo la request con la instancia de RegisterReservationDTO
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	//Valido los datos
	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.RegisterReservation(ctx, params.UserID, params.Name, params.Password, params.Email, params.TableNumber, params.ReservationDate)
	if err != nil {
		if err == service.ErrTableNotAvailable {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}

		log.Println("Error during reservation registration:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Reservation registered successfully"})
}

func (a *API) RemoveReservation(c echo.Context) error {

	ctx := c.Request().Context()

	params := dtos.RemoveReservationDTO{}

	//Linkeo la request con la instancia de RemoveReservationDTO
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.RemoveReservation(ctx, params.UserID)
	if err != nil {
		if err == service.ErrReservationNotFound {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing registration:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Reservation removed successfully"})
}

func (a *API) AddTable(c echo.Context) error {

	ctx := c.Request().Context()

	params := dtos.AddTableDTO{}

	//Linkeo la request con la instancia de RegisterReservationDTO
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	//Valido los datos
	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.AddTable(ctx, params.Number, params.Seats, params.Location)
	if err != nil {
		if err == service.ErrAddingTable {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}

		log.Println("Error while adding a table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Table added successfully"})
}

func (a *API) RemoveTable(c echo.Context) error {

	ctx := c.Request().Context()

	params := dtos.RemoveTableDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}
	err = a.serv.RemoveTable(ctx, params.Number)
	if err != nil {
		if err == service.ErrRemovingTable {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Table removed successfully"})

}

func (a *API) RemoveUser(c echo.Context) error {

	ctx := c.Request().Context()

	params := dtos.RemoveUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.RemoveUser(ctx, params.UserID)
	if err != nil {
		if err == service.ErrReservationNotFound {
			return c.JSON(http.StatusConflict, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing user:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "User removed successfully"})
}
