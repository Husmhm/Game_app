package config

import "time"

const (
	Jwt_SignKey                = "jwt_secret"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AuthMiddleWareConrexKey    = "claims"
)
