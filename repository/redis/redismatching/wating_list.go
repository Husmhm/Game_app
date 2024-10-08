package redismatching

import (
	"context"
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/richerror"
	"github.com/redis/go-redis/v9"
	"time"
)

// TODO - add to config in usecase layer...
const WatingListPrefix = "watinglist"

func (d DB) AddToWatingList(userID uint, category entity.Category) error {
	const op = "redismatching.AddToWatingList"

	_, err := d.adapter.Client().ZAdd(context.Background(), fmt.Sprintf("%s:%s", WatingListPrefix, category), redis.Z{
		Score:  float64(time.Now().UnixMicro()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}
