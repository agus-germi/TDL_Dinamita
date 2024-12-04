package jwtutils

import (
	"errors"
	"time"

	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/agus-germi/TDL_Dinamita/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidKey          = jwt.ErrInvalidKey
	ErrParsingToken        = errors.New("error parsing token")
	ErrInvalidClaimsFormat = errors.New("invalid claims format")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("token has expired")
	ErrMissingToken        = errors.New("missing token")
)

var (
	jwtSecretKey   string
	expirationTime int64
)

// In Go, the init() function is a special feature designed to perform initial
// configurations of the package, such as initializing variables, validations,
// or necessary connections before using the rest of the code in the package.
// The init() function is called automatically when the package is imported
// into another package.
func init() {
	logger.Log.Debug("Executing init() function of 'jwtuilts' package: Loading JWT Secret Key and Expiration Time from '.env' file")

	jwtSecretKey, err := utils.GetEnv("JWT_SECRET_KEY")
	if err != nil || jwtSecretKey == "" {
		logger.Log.Fatalf("JWT_SECRET_KEY is not set or invalid: %v", err)
	}
	logger.Log.Debugf("Value read from JWT_SECRET_KEY: %s", jwtSecretKey)

	timeToAddStr, err := utils.GetEnv("JWT_EXPIRATION_TIME")
	if err != nil || timeToAddStr == "" {
		logger.Log.Fatalf("JWT_EXPIRATION_TIME is not set or invalid: %v", err)
	}
	logger.Log.Debugf("Value read from JWT_EXPIRATION_TIME: %s", timeToAddStr)

	timeToAdd, err := time.ParseDuration(timeToAddStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to parse duration of JWT_EXPIRATION_TIME environment variable: %v", err)
	}

	expirationTime = time.Now().Add(timeToAdd).Unix() // Expiration time according Unix format.

	logger.Log.Infof("JWT Secret Key and Expiration Time loaded successfully from '.env' file.")
}

func SignedLoginToken(u *models.User) (string, error) {
	//HS256 > viable porque el servidor que creo que el certificado tambien lo validará

	//Claims: estructura de datos que se puede firmar y validar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"role_id": u.RoleID,
		"name":    u.Name,
		"exp":     expirationTime,
	})
	jwt, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	logger.Log.Debugf("Session token created:", jwt)
	return jwt, nil
}

func GetClaimsFromToken(c echo.Context) (jwt.MapClaims, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, ErrMissingToken
	}

	logger.Log.Debugf("Session token with 'Bearer' at the beginning:", token)
	// Deleting "Bearer " from the beginning of the token
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	logger.Log.Debugf("Session token without 'Bearer':", token)

	claims, err := parseLoginJWT(token)
	if err != nil {
		logger.Log.Errorf("Error while parsing token:", err)
		return nil, err
	}

	return claims, nil
}

func parseLoginJWT(value string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	logger.Log.Debugf("Parsed token:", token)

	if err != nil {
		// Distinction between expired token and other errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	// Check if the claims are valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaimsFormat
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
