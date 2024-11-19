package main

import (
	"gameApp/adapter/redis"
	"gameApp/config"
	"gameApp/delivery/grpcserver/presenceserver"
	"gameApp/repository/redis/redispresence"
	"gameApp/service/presenceservice"
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	server := presenceserver.New(presenceSvc)
	server.Start()
}
