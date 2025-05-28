package utils

import (
	"os"

	"github.com/OleksandrBob/nextseasonlist/shared/token"
	"github.com/dgrijalva/jwt-go"
)

var accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return token.ValidateToken(tokenString, accessTokenSecret)
}
