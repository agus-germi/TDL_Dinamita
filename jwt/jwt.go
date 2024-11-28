package jwt

import (
	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const key = "01234567890123456789012345678901" //TODO:  chequear si hay una mejor forma o una mejor key.

// SignedLoginToken genera un token firmado con el email y el nombre del usuario
func SignedLoginToken(u *models.User) (string, error) {
	//HS256 > viable porque el servidor que creo que el certificado tambien lo validará

	//Claims: estructura de datos que se puede firmar y validar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": u.Email,
		"name":  u.Name,
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
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil

}
