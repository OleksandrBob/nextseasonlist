package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
var refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

func GenerateAccessToken(userID string) (string, error) {
	calims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	}
	return generateToken(calims, accessTokenSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	calims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	}
	return generateToken(calims, refreshTokenSecret)
}

func generateToken(claims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return validateToken(tokenString, accessTokenSecret)
}

func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return validateToken(tokenString, refreshTokenSecret)
}

func validateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
