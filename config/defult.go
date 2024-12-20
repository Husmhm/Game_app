package config

import "time"

var defaultConfig = map[string]interface{}{
	"auth.access_subject":                    AccessTokenSubject,
	"auth.refresh_subject":                   RefreshTokenSubject,
	"auth.access_expiration_time":            AccessTokenExpireDuration,
	"auth.refresh_expiration_time":           RefreshTokenExpireDuration,
	"application.gracefull_shutdown_timeout": time.Second * 5,
}
