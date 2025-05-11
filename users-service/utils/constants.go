package utils

import "time"

const (
	AccessTokenDurationTime  = 10 * time.Minute
	RefreshTokenDurationTime = (7 * 24 * time.Hour) // 1 week
)

const (
	AccessTokenName  = "accessToken"
	RefreshTokenName = "refreshToken"
)

const (
	UserIdClaim     = "sub"
	ExpirationClaim = "exp"
)

const (
	JobMaxDurationTime = 1 * time.Minute
	JobRecuuringPeriod = 24 * time.Hour
)
