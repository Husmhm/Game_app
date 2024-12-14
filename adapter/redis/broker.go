package redis

import (
	"context"
	"gameApp/entity"
	"github.com/labstack/gommon/log"
	"time"
)

func (a Adapter) Publish(event entity.Event, payload string) {
	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancle()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf("publish err: %v\n", err)
		// TODO - log
		//TODO -update metrics
	}
	//TODO -update metrics
}
