package jwt

import (
	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const key = "01234567890123456789012345678901" // chequear si hay una mejor forma o una mejor key.

func SignedLoginToken(usr *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": usr.Email,
		"name":  usr.Name,
	})

	return token.SignedString([]byte(key))
}
