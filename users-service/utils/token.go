package utils

import (
	"os"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shared/token"
	"github.com/dgrijalva/jwt-go"
)

var accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
var refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

func GenerateAccessToken(userID string, userRoles []string) (string, error) {
	calims := jwt.MapClaims{
		UserIdClaim:     userID,
		RolesClaim:      userRoles,
		ExpirationClaim: time.Now().Add(AccessTokenDurationTime).Unix(),
	}

	return generateToken(calims, accessTokenSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	calims := jwt.MapClaims{
		UserIdClaim:     userID,
		ExpirationClaim: time.Now().Add(RefreshTokenDurationTime).Unix(),
	}

	return generateToken(calims, refreshTokenSecret)
}

func generateToken(claims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return token.ValidateToken(tokenString, refreshTokenSecret)
}
