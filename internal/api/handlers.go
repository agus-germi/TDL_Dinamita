package api

import (
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

var errorResponses = map[error]int{
	service.ErrInvalidCredentials:   http.StatusBadRequest,
	service.ErrUserNotFound:         http.StatusNotFound,
	service.ErrUserAlreadyExists:    http.StatusConflict,
	service.ErrReservationNotFound:  http.StatusNotFound,
	service.ErrTableAlreadyExists:   http.StatusConflict,
	service.ErrTableNotAvailable:    http.StatusConflict,
	service.ErrTableNotFound:        http.StatusNotFound,
	service.ErrUserRoleAlreadyAdded: http.StatusConflict,
	service.ErrInvalidPermission:    http.StatusForbidden,
}

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
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "User registered successfully"})
}

func (a *API) CreateReservation(c echo.Context) error {

	ctx := c.Request().Context() // obtengo el contexto del objeto Request que viene con la petición HTTP

	params := dtos.CreateReservationDTO{}

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

	// Get userID from cookie
	claims, err = jwtutils.GetClaimsFromCookie(c)
	if err != nil {
		log.Println("Error while getting the user ID from the cookie:", err)
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Unauthorized"})
	}
	userID := claims["user_id"].(int64)

	err = a.serv.RegisterReservation(ctx, userID, params.TableNumber, params.ReservationDate)
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

	email := claims["email"].(string)
	err = notification.SendReservationConfirmationEmail(email, emailBody)
	if err != nil {
		log.Println("Failed to send confirmation email:", err)
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

	err = a.serv.RemoveReservation(ctx, params.ID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing registration:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Reservation removed successfully"})
}

func (a *API) AddTable(c echo.Context) error {
	email, err := getEmailFromCookie(c)
	if err != nil {
		log.Println("Error while getting the email from the cookie:", err)
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Unauthorized"})
	}
	println("email", email)

	ctx := c.Request().Context()
	params := dtos.AddTableDTO{}

	//Linkeo la request con la instancia de RegisterReservationDTO
	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	//Valido los datos
	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.AddTable(ctx, params.Number, params.Seats, params.Location, email)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while adding a table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Table added successfully"})
}

func (a *API) RemoveTable(c echo.Context) error {
	email, err := getEmailFromCookie(c)
	if err != nil {
		log.Println("Error while getting the email from the cookie:", err)
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Unauthorized"})
	}

	ctx := c.Request().Context()
	params := dtos.RemoveTableDTO{}

	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.RemoveTable(ctx, params.Number, email)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing table:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Table removed successfully"})
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
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while removing user:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "User removed successfully"})
}

func (a *API) AddUserRole(c echo.Context) error {

	email, err := getEmailFromCookie(c)
	if err != nil {
		log.Println("Error while getting the email from the cookie:", err)
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Unauthorized"})
	}

	ctx := c.Request().Context()
	params := dtos.UserRoleDTO{}

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

	err = a.serv.UpdateUserRole(ctx, params.UserID, params.RoleID, email)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error while assigning a role to the user:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Role assigned to the user successfully"})
}

func (a *API) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()

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

	usr, err := a.serv.LoginUser(ctx, params.Email, params.Password)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			log.Println("Error trying to login:", err)
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		log.Println("Error trying to login:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	token, err := jwt.SignedLoginToken(usr)
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

func (a *API) GetAllReservationsOfUser(c echo.Context) error {
	ctx := c.Request().Context()

	base := 10
	bitSize := 64

	userID, err := strconv.ParseInt(c.Param("id"), base, bitSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID"})
	}

	params := dtos.GetReservationsDTO{UserID: userID}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	reservations, err := a.serv.GetReservationsByUserID(ctx, params.UserID)
	if err != nil {
		if statusCode, ok := errorResponses[err]; ok {
			return c.JSON(statusCode, responseMessage{Message: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	dtoReservations := convertReservationsToDTO(reservations)

	return c.JSON(http.StatusOK, dtoReservations)
}

func convertReservationsToDTO(reservations *[]models.Reservation) *[]dtos.ReservationDTO {
	if reservations == nil {
		return &[]dtos.ReservationDTO{}
	}

	dtoReservations := make([]dtos.ReservationDTO, len(*reservations))
	for i, reservation := range *reservations {
		dtoReservations[i] = dtos.ReservationDTO{
			ID:              reservation.ID,
			UserID:          reservation.UserID,
			TableNumber:     reservation.TableNumber,
			ReservationDate: reservation.ReservationDate,
		}
	}
	return &dtoReservations
}

// funciones aux

// getEmailFromCookie obtiene el email del usuario a partir de la cookie de autenticación
func getEmailFromCookie(c echo.Context) (string, error) {

	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return "", err
	}
	claims, err := jwt.ParseLoginJWT(cookie.Value)
	if err != nil {
		return "", err
	}
	email := claims["email"].(string)
	return email, nil
}
