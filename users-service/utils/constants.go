package utils

import "time"

var (
	AccessTokenDurationTime  = 10 * time.Minute
	RefreshTokenDurationTime = (7 * 24 * time.Hour) // 1 week
)

var (
	AccessTokenName  = "accessToken"
	RefreshTokenName = "refreshToken"
)

var (
	UserIdClaim     = "user_id"
	ExpirationClaim = "exp"
)
