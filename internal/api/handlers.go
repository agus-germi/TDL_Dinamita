package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.RegisterUser(ctx, params.Name, params.Password, params.Email)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while registering user: %v")
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "User registered successfully"})
}

func (a *API) LoginUser(c echo.Context) error {
	params := dtos.LoginUserDTO{}

	err := c.Bind(&params)
	if err != nil {
		a.log.Errorf("[User Login] error while binding info of DTO: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		a.log.Errorf("[User Login] error during data validation: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	usr, err := a.serv.LoginUser(ctx, params.Email, params.Password)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error trying to login: %v")
	}

	token, err := jwtutils.SignedLoginToken(usr)
	if err != nil {
		a.log.Errorf("Error trying to create a jwt:", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
	}

	a.log.Infof("[User Login] User logged successfully: %v", usr.Email)
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User logged successfully", "token": token})
}

func (a *API) DeleteUser(c echo.Context) error {

	clientUserID, err := getContextValueAsInt64(c, "user_id", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	clientRoleID, err := getContextValueAsInt64(c, "role", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	userIDToDelete, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Delete User] Failed to parse user ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID"})
	}

	if clientRoleID != adminRoleID && clientUserID != userIDToDelete {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can only delete your own account"})
	}

	ctx := c.Request().Context()
	err = a.serv.RemoveUser(ctx, userIDToDelete)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while deleting user: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "User deleted successfully"})
}

func (a *API) UpdateUserRole(c echo.Context) error {

	userIDToUpdate, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Update User] Failed to parse user ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID"})
	}

	ctx := c.Request().Context()
	params := dtos.UpdateUserRoleDTO{}

	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	err = a.serv.UpdateUserRole(ctx, userIDToUpdate, params.NewRoleID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while assigning a new role to the user: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "New role assigned to the user successfully"})
}

func (a *API) GetReservationsOfUser(c echo.Context) error {

	clientUserID, err := getContextValueAsInt64(c, "user_id", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	clientRoleID, err := getContextValueAsInt64(c, "role", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	userID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Get Reservations of User] Failed to parse user ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid user ID"})
	}

	if clientRoleID != adminRoleID && clientUserID != userID {
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Permission denied: you can only see your own reservations"})
	}

	ctx := c.Request().Context()
	reservations, err := a.serv.GetReservationsByUserID(ctx, userID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while getting reservations: %v")
	}

	dtoReservations := convertReservationsToDTO(reservations)

	return c.JSON(http.StatusOK, dtoReservations)
}

// Reservation endpoints
func (a *API) CreateReservation(c echo.Context) error {

	clientUserID, err := getContextValueAsInt64(c, "user_id", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}
	clientEmail, ok := c.Get("email").(string)
	a.log.Debugf("[Create Reservation] Client Email:", clientEmail)
	if !ok {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: "Invalid client email in context"})
	}

	params := dtos.CreateReservationDTO{}

	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		a.log.Errorf("[Create Reservation] error during data validation: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	reservationDate, _ := time.Parse(time.RFC3339, params.ReservationDate) // Think if we need to specify the time zone (-03:00 for Buenos Aires)
	a.log.Debugf("[Create Reservation] Reservation Date:", reservationDate)

	ctx := c.Request().Context()
	start := time.Now()
	promotionID := int(params.PromotionID)
	err = a.serv.MakeReservation(ctx, clientUserID, params.TableNumber, reservationDate, promotionID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error during reservation registration: %v")
	}
	duration := time.Since(start)
	a.log.Infof("[Create Reservation] Reservation registered in %v", duration)

	emailBody := fmt.Sprintf(
		"Hello!<br><br>Your reservation for table %d on %s has been confirmed.<br>Thank you!",
		params.TableNumber, reservationDate.Format("2006-01-02 15:04"),
	)

	err = notification.SendReservationConfirmationEmail(clientEmail, emailBody)
	if err != nil {
		a.log.Errorf("Failed to send confirmation email:", err)
	}
	return c.JSON(http.StatusCreated, responseMessage{Message: "Reservation registered successfully"})
}

func (a *API) DeleteReservation(c echo.Context) error {

	clientUserID, err := getContextValueAsInt64(c, "user_id", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	clientRoleID, err := getContextValueAsInt64(c, "role", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	reservationID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Delete Reservation] Failed to parse reservation ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid reservation ID"})
	}

	ctx := c.Request().Context()
	err = a.checkUserPermissionToCancelReservation(ctx, clientUserID, clientRoleID, reservationID)
	if err != nil {
		return c.JSON(http.StatusForbidden, responseMessage{Message: err.Error()})
	}

	err = a.serv.CancelReservation(ctx, reservationID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while canceling reservation: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Reservation canceled successfully"})
}

// Table endpoints
func (a *API) CreateTable(c echo.Context) error {

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
	err = a.serv.AddTable(ctx, params.Number, params.Seats, params.Description)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while adding a table: %v")
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Table added successfully"})
}

func (a *API) DeleteTable(c echo.Context) error {

	tableID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Delete Table] Failed to parse table ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid table ID"})
	}

	ctx := c.Request().Context()
	err = a.serv.RemoveTable(ctx, tableID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while deleting table: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Table deleted successfully"})
}

// Menu endpoints
func (a *API) AddDishToMenu(c echo.Context) error {

	params := dtos.DishDTO{}

	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.AddDishToMenu(ctx, params.Name, params.Price, params.Description)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while adding a dish to the menu: %v")
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Dish added to the menu successfully"})

}

func (a *API) RemoveDishFromMenu(c echo.Context) error {

	dishID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Delete Dish] Failed to parse dish ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid dish ID"})
	}

	ctx := c.Request().Context()
	err = a.serv.RemoveDish(ctx, dishID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while deleting dish: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Dish deleted successfully"})
}

func (a *API) UpdateDishInMenu(c echo.Context) error {

	dishID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Update Dish In Menu] Failed to parse dish ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid dish ID"})
	}

	ctx := c.Request().Context()
	params := dtos.DishDTO{} // uso el mismo DTO porque considero que se pueden modificar todos los campos del plato

	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}
	err = a.serv.UpdateDish(ctx, dishID, params.Name, params.Price, params.Description)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while updating dish: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Dish updated successfully"})

}

func (a *API) GetDishesInMenu(c echo.Context) error {
	ctx := c.Request().Context()
	dishes, err := a.serv.GetDishes(ctx)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while getting dishes: %v")
	}
	//convertir dishes a DTO
	dtoDishes := convertDishesToDTO(dishes)
	return c.JSON(http.StatusOK, dtoDishes)

}

func (a *API) GetTables(c echo.Context) error {
	ctx := c.Request().Context()
	tables, err := a.serv.GetAvailableTables(ctx)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while fetching available tables: %v")
	}

	return c.JSON(http.StatusOK, tables)
}

// Time slots endpoint
func (a *API) GetTimeSlots(c echo.Context) error {
	ctx := c.Request().Context()

	timeSlots, err := a.serv.GetTimeSlots(ctx)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while getting time slots: %v")
	}

	dtoTimeSlots := convertTimeSlotsToDTO(timeSlots)

	return c.JSON(http.StatusOK, dtoTimeSlots)
}

// Opinions endpoints
func (a *API) CreateOpinion(c echo.Context) error {

	clientUserID, err := getContextValueAsInt64(c, "user_id", c.Logger())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responseMessage{Message: err.Error()})
	}

	params := dtos.CreateOpinionDTO{}
	err = c.Bind(&params)
	if err != nil {
		a.log.Errorf("[Create Opinion] Error parsing request body: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		a.log.Errorf("[Create Opinion] Validation error: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.CreateOpinion(ctx, clientUserID, params.Opinion, params.Rating)
	if err != nil {
		a.log.Errorf("[Create Opinion] Error while creating opinion: %v", err)
		return a.handleErrorFromService(c, err, "Error while creating opinion: %v")
	}

	// Responder al cliente
	return c.JSON(http.StatusCreated, responseMessage{Message: "Opinion created successfully"})
}

func (a *API) GetOpinions(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch all opinions (you could add filters here if needed)
	opinions, err := a.serv.GetOpinions(ctx)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while fetching opinions: %v")
	}

	return c.JSON(http.StatusOK, opinions)
}

// Promotions endpoint
func (a *API) CreatePromotion(c echo.Context) error {

	params := dtos.CreatePromotionDTO{}
	err := c.Bind(&params)
	if err != nil {
		a.log.Errorf("[Create Promotion] Error parsing request body: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid request"})
	}

	err = a.dataValidator.Struct(params)
	if err != nil {
		a.log.Errorf("[Create Promotion] Validation error: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.serv.CreatePromotion(ctx, params.Description, params.StartDate, params.DueDate, params.Discount)
	if err != nil {
		a.log.Errorf("[Create Promotion] Error while creating promotion: %v", err)
		return a.handleErrorFromService(c, err, "Error while creating promotion: %v")
	}

	return c.JSON(http.StatusCreated, responseMessage{Message: "Promotion created successfully"})
}

func (a *API) DeletePromotion(c echo.Context) error {

	promotionID, err := parseID(c, "id")
	if err != nil {
		a.log.Errorf("[Delete Promotion] Failed to parse dish ID: %v", err)
		return c.JSON(http.StatusBadRequest, responseMessage{Message: "Invalid dish ID"})
	}

	ctx := c.Request().Context()
	err = a.serv.DeletePromotion(ctx, promotionID)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while deleting Promotion: %v")
	}

	return c.JSON(http.StatusOK, responseMessage{Message: "Promotion Deleted successfully"})
}

func (a *API) GetPromotions(c echo.Context) error {
	ctx := c.Request().Context()

	opinions, err := a.serv.GetPromotions(ctx)
	if err != nil {
		return a.handleErrorFromService(c, err, "Error while fetching promotions: %v")
	}

	return c.JSON(http.StatusOK, opinions)
}

// Auxiliary functions
func (a *API) handleErrorFromService(c echo.Context, err error, debugMsg string) error {
	if statusCode, ok := errorResponses[err]; ok {
		return c.JSON(statusCode, responseMessage{Message: err.Error()})
	}

	a.log.Errorf(debugMsg, err)
	return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Internal server error"})
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
			Promotion:       reservation.Promotion,
		}
	}
	return &dtoReservations
}

func convertDishesToDTO(dishes *[]models.Dish) *[]dtos.DishDTO {
	if dishes == nil {
		return &[]dtos.DishDTO{}
	}

	dtoDishes := make([]dtos.DishDTO, len(*dishes))
	for i, dish := range *dishes {
		dtoDishes[i] = dtos.DishDTO{
			Name:        dish.Name,
			Price:       dish.Price,
			Description: dish.Description,
		}
	}
	return &dtoDishes
}

func convertTimeSlotsToDTO(timeSlots *[]models.TimeSlot) []map[string]interface{} {
	dto := make([]map[string]interface{}, len(*timeSlots))
	for i, ts := range *timeSlots {
		dto[i] = map[string]interface{}{
			"id":   ts.ID,
			"time": ts.Time,
		}
	}
	return dto
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
		return errors.New("you are not allowed to cancel this reservation")
	}

	return nil
}

func parseID(c echo.Context, paramName string) (int64, error) {
	// Parses the ID from the URI and returns it as an int64
	base := 10
	bitSize := 64
	idStr := c.Param(paramName)
	id, err := strconv.ParseInt(idStr, base, bitSize)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", paramName, err)
	}
	return id, nil
}

func getContextValueAsInt64(c echo.Context, key string, log echo.Logger) (int64, error) {
	value, ok := c.Get(key).(float64) //Asercion de tipo
	log.Debugf("[%s] Retrieved value: %v", key, value)

	if !ok {
		return 0, fmt.Errorf("invalid %s in context", key)
	}

	valueInt64 := int64(value)
	log.Debugf("[%s] Converted value: %d", key, valueInt64)
	return valueInt64, nil
}
