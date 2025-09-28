package utils

import (
	"os"
	"time"

	"github.com/OleksandrBob/nextseasonlist/shared/token"
	"github.com/dgrijalva/jwt-go"
)

func GenerateAccessToken(userID string, userRoles []string) (string, error) {
	calims := jwt.MapClaims{
		UserIdClaim:     userID,
		RolesClaim:      userRoles,
		ExpirationClaim: time.Now().UTC().Add(AccessTokenDurationTime).Unix(),
	}

	return generateToken(calims, []byte(os.Getenv("ACCESS_TOKEN_SECRET")))
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		UserIdClaim:     userID,
		ExpirationClaim: time.Now().UTC().Add(RefreshTokenDurationTime).Unix(),
	}

	return generateToken(claims, []byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}

func generateToken(claims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return token.ValidateToken(tokenString, []byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}
