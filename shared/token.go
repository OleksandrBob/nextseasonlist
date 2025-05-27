package shared

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func validateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return secret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return claims, errors.New("token expired")
			}
			return nil, errors.New("invalid token: " + err.Error())
		}
		return nil, err
	}

	return claims, nil
}
