package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/agus-germi/TDL_Dinamita/internal/api/dtos"
	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/agus-germi/TDL_Dinamita/internal/service/notification"
	"github.com/agus-germi/TDL_Dinamita/jwtutils"
	"github.com/labstack/echo/v4"
)

const adminRoleID = 1

var errorResponses = map[error]int{
	service.ErrInvalidCredentials:   http.StatusBadRequest,
	service.ErrUserNotFound:         http.StatusNotFound,
	service.ErrUserAlreadyExists:    http.StatusConflict,
	service.ErrReservationNotFound:  http.StatusNotFound,
	service.ErrTableAlreadyExists:   http.StatusConflict,
	service.ErrTableNotAvailable:    http.StatusConflict,
	service.ErrTableNotFound:        http.StatusNotFound,
	service.ErrUserRoleAlreadyAdded: http.StatusConflict,
}

type responseMessage struct {
	Message string `json:"message"`
}

// User endpoints
func (a *API) RegisterUser(c echo.Context) error {
	params := dtos.RegisterUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"}) // otra opcion--> Message: err.Error()
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.RegisterUser(ctx, params.Name, params.Password, params.Email)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "User registered successfully"})
}

func (a *API) LoginUser(c echo.Context) error {
	params := dtos.LoginUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	usr, err := a.serv.LoginUser(ctx, params.Email, params.Password)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			log.Println("Error trying to login:", err)
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error trying to login:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	token, err := jwtutils.SignedLoginToken(usr)
	if err != nil {
		log.Println("Error trying to create a jwt:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		SameSite: http.SameSiteNoneMode, // Indica que la cookie no debe ser enviada en una petición de un sitio diferente al que la generó.
		Secure:   true,                  // Indica que la cookie solo debe ser enviada al servidor (nuestra API) si la conexión se realiza a través de HTTPS.
		HttpOnly: true,                  // Previene que la cookie sea accesible desde JavaScript ejecutado en el navegador. (impide que scripts maliciosos lean o manipulen las cookies.)
		Path:     "/",                   // Hacemos accesible la cookie para todos los endpointsde la aplicacion.
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}

// Este endpoint se llamaria cuando un usuario tomase la decision de eliminar su cuenta.
func (a *API) DeleteUser(c echo.Context) error {
	clientUserID, ok := c.Get("user_id").(int64) // Aserción de tipo
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client user ID in context"})
	}

	clientRoleID, ok := c.Get("role").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
	}

	log.Println("[Delete Reservation] Client User ID:", clientUserID)
	log.Println("[Delete Reservation] Client Role ID:", clientRoleID)

	base := 10
	bitSize := 64

	userIDToDelete, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID to delete"})
	}

	if clientRoleID != adminRoleID && clientUserID != userIDToDelete {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can only delete your own account"})
	}

	ctx := c.Request().Context()
	err = a.serv.RemoveUser(ctx, userIDToDelete)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while deleting user:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "User deleted successfully"})
}

func (a *API) UpdateUserRole(c echo.Context) error {
	clientRoleID, ok := c.Get("role").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
	}

	log.Println("[Update User Role] Client Role ID:", clientRoleID)

	if clientRoleID != adminRoleID {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can't assign a new role to a user"})
	}

	base := 10
	bitSize := 64

	userIDToUpdate, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID to update"})
	}

	ctx := c.Request().Context()
	params := dtos.UpdateUserRoleDTO{}

	// Linkeo la request con la instancia de UserRoleDTO
	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	//Valido los datos
	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.UpdateUserRole(ctx, userIDToUpdate, params.NewRoleID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while assigning a new role to the user:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "New role assigned to the user successfully"})
}

func (a *API) GetReservationsOfUser(c echo.Context) error {
	clientUserID, ok := c.Get("user_id").(int64) // Aserción de tipo
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client user ID in context"})
	}

	clientRoleID, ok := c.Get("user_id").(int64) // Aserción de tipo
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client user ID in context"})
	}

	log.Println("[Get Reservations of User] Client User ID:", clientUserID)
	log.Println("[Get Reservations of User] Client Role ID:", clientRoleID)

	base := 10
	bitSize := 64

	userID, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID"})
	}

	if clientRoleID != adminRoleID && clientUserID != userID {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can only see your own reservations"})
	}

	ctx := c.Request().Context()
	reservations, err := a.serv.GetReservationsByUserID(ctx, userID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	dtoReservations := convertReservationsToDTO(reservations)

	return c.JSON(http.StatusOK, dtoReservations)
}

// Reservation endpoints
func (a *API) CreateReservation(c echo.Context) error {
	clientUserID, ok := c.Get("user_id").(int64) // Aserción de tipo
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client user ID in context"})
	}

	clientEmail, ok := c.Get("email").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client email in context"})
	}

	log.Println("[Create Reservation] User ID:", clientUserID)
	log.Println("[Create Reservation] Email:", clientEmail)

	params := dtos.CreateReservationDTO{}

	//Linkeo la request con la instancia de CreateReservationDTO
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context() // obtengo el contexto del objeto Request que viene con la petición HTTP
	err = a.serv.MakeReservation(ctx, clientUserID, params.TableNumber, params.ReservationDate)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error during reservation registration:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	emailBody := fmt.Sprintf(
		"Hello!<br><br>Your reservation for table %d on %s has been confirmed.<br>Thank you!",
		params.TableNumber, params.ReservationDate.Format("2006-01-02 15:04"),
	)

	err = notification.SendReservationConfirmationEmail(clientEmail, emailBody)
	if err != nil {
		log.Println("Failed to send confirmation email:", err)
	}
	return c.JSON(http.StatusCreated, responseMessage{Message: "Reservation registered successfully"})
}

func (a *API) DeleteReservation(c echo.Context) error {
	clientUserID, ok := c.Get("user_id").(int64) // Aserción de tipo
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client user ID in context"})
	}

	clientRoleID, ok := c.Get("role").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
	}

	log.Println("[Delete Reservation] User ID:", clientUserID)
	log.Println("[Delete Reservation] Role ID:", clientRoleID)

	base := 10
	bitSize := 64

	reservationID, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid reservation ID"})
	}

	ctx := c.Request().Context()
	err = a.checkUserPermissionToCancelReservation(ctx, clientUserID, clientRoleID, reservationID)
	if err != nil {
		return c.JSON(http.StatusForbidden, responseMessage{Message: err.Error()})
	}

	err = a.serv.CancelReservation(ctx, reservationID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while cancelling reservation:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Reservation canceled successfully"})
}

// Table endpoints
func (a *API) CreateTable(c echo.Context) error {
	clientRoleID, ok := c.Get("role").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
	}

	log.Println("[Create Table] Role ID:", clientRoleID)

	if clientRoleID != adminRoleID {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can't add a new table"})
	}

	params := dtos.CreateTableDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.AddTable(ctx, params.Number, params.Seats, params.Location)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while adding a table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Table added successfully"})
}

func (a *API) DeleteTable(c echo.Context) error {
	clientRoleID, ok := c.Get("role").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client role ID in context"})
	}

	log.Println("[Delete Table] Role ID:", clientRoleID)

	if clientRoleID != adminRoleID {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can't add a new table"})
	}

	base := 10
	bitSize := 64

	tableID, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid table ID"})
	}

	ctx := c.Request().Context()
	err = a.serv.RemoveTable(ctx, tableID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while deleting table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Table deleted successfully"})
}

func convertReservationsToDTO(reservations *[]models.Reservation) *[]dtos.ReservationDTO {
	if reservations == nil {
		return &[]dtos.ReservationDTO{}
	}

	dtoReservations := make([]dtos.ReservationDTO, len(*reservations))
	for i, reservation := range *reservations {
		dtoReservations[i] = dtos.ReservationDTO{
			ID:              reservation.ID,
			TableNumber:     reservation.TableNumber,
			ReservationDate: reservation.ReservationDate,
		}
	}
	return &dtoReservations
}

func (a *API) checkUserPermissionToCancelReservation(ctx context.Context, userID, roleID, reservationID int64) error {
	// An admin user can cancel any reservation
	if roleID == adminRoleID {
		return nil
	}

	// A regular user can only cancel their own reservations
	reservation, err := a.serv.GetReservationByID(ctx, reservationID)
	if err != nil {
		return err
	}

	// Verfify if the user is the owner of the reservation
	if reservation.UserID != userID {
		return errors.New("You are not allowed to cancel this reservation")
	}

	return nil
}
