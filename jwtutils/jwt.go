package jwt

import (
	"errors"
	"time"

	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const key = "01234567890123456789012345678901" // TODO: Cambiar por una key más segura, idealmente almacenada en una variable de entorno

var (
	ErrInvalidKey          = jwt.ErrInvalidKey
	ErrParsingToken        = errors.New("error parsing token")
	ErrInvalidClaimsFormat = errors.New("invalid claims format")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("token has expired")
)

// SignedLoginToken genera un token firmado con el id, email y el nombre del usuario
func SignedLoginToken(u *models.User) (string, error) {
	//HS256 > viable porque el servidor que creo que el certificado tambien lo validará

	// Expiration time: 24 hours
	expirationTime := time.Now().Add(24 * time.Hour).Unix() // Expiration time according Unix format.

	//Claims: estructura de datos que se puede firmar y validar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"name":    u.Name,
		"exp":     expirationTime,
	})
	jwt, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	println("token", jwt)
	return jwt, nil
}

func ParseLoginJWT(value string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
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

func GetClaimsFromCookie(c echo.Context) (jwt.MapClaims, error) {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return nil, err
	}
	claims, err := ParseLoginJWT(cookie.Value)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
