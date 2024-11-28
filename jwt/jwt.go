package jwt

import (
	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const key = "01234567890123456789012345678901" //TODO:  chequear si hay una mejor forma o una mejor key.

func SignedLoginToken(usr *models.User) (string, error) {

	//HS256 > viable porque el servidor que creo que el certificado tambien lo validará
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": usr.Email,
		"name":  usr.Name,
	})

	jwtString, err := token.SignedString([]byte(key))
	if err != nil {
		return " ", err
	}
	return jwtString, nil
}
