package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"gameApp/adapter/redis"
	"gameApp/config"
	"gameApp/contract/protogolang/matching"
	"gameApp/entity"
	"gameApp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := "matching.users_matched"
	mu := entity.MatchedUsers{
		Category: entity.FootballCategory,
		UserIDs:  []uint{1, 4},
	}

	pbmu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
	}
	payload, err := proto.Marshal(&pbmu)
	if err != nil {
		panic(err)
	}

	payloadStr := base64.StdEncoding.EncodeToString(payload)
	fmt.Println(payloadStr)
	if err := redisAdapter.Client().Publish(context.Background(), topic, payloadStr).Err(); err != nil {
		panic(fmt.Sprintf("publisn err: %v", err))
	}
}
