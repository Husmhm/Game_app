package redismatching

import (
	"context"
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/richerror"
	"gameApp/pkg/timestamp"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// TODO - add to config in usecase layer...
const WatingListPrefix = "watinglist"

func (d DB) AddToWatingList(userID uint, category entity.Category) error {
	const op = richerror.Op("redismatching.AddToWatingList")

	_, err := d.adapter.Client().ZAdd(context.Background(), getCtegory(category), redis.Z{
		Score:  float64(timestamp.Now()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = richerror.Op("redismatching.GetWaitingListByCategory")
	//d.adapter.Client().ZRangeWithScores()
	min := strconv.Itoa(int(timestamp.Add(-2 * time.Hour)))
	max := strconv.Itoa(int(timestamp.Now()))
	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx, getCtegory(category), &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	var result []entity.WaitingMember

	for _, l := range list {
		userID, _ := strconv.Atoi(l.Member.(string))
		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			TimeStamp: int64(l.Score),
			Category:  category,
		})
	}
	return result, nil

}

func getCtegory(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WatingListPrefix, category)
}

func (d DB) RemoveUsersFromWaitingList(category entity.Category, userIDs []uint) {
	// TODO - add 5 to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	members := make([]any, 0)
	for _, userID := range userIDs {
		members = append(members, strconv.Itoa(int(userID)))
	}

	numberOfRemovedMembers, err := d.adapter.Client().ZRem(ctx, getCtegory(category), members...).Result()
	if err != nil {
		log.Errorf("remve from waiting list error: %v \n", err)
		// TODO - update metrics
	}

	log.Printf("%d items removed from %s", numberOfRemovedMembers, getCtegory(category))
	// TODO - Updete metrics
}
