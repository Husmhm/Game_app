package config

import (
	"gameApp/adapter/redis"
	"gameApp/repository/mysql"
	"gameApp/scheduler"
	"gameApp/service/authservice"
	"gameApp/service/matchingservice"
	"gameApp/service/presenceservice"
	"time"
)

type Application struct {
	GraceFullShutdownTimeout time.Duration `koanf:"gracefull_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}
type Config struct {
	Application     Application            `koanf:"application"`
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
	PresenceService presenceservice.Config `koanf:"presence_service"`
	Schduller       scheduler.Config       `koanf:"schduller"`
}
