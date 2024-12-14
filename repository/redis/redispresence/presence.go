package redispresence

import (
	"context"
	"gameApp/pkg/richerror"
	"gameApp/pkg/timestamp"
	"time"
)

func (d DB) Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error {
	const op = richerror.Op("redispresence.Upsert")
	_, err := d.adapter.Client().Set(ctx, key, timestamp, expTime).Result()
	if err != nil {
		richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

func (d DB) GetPresence(ctx context.Context, prefixkey string, userIDs []uint) (map[uint]int64, error) {
	const op = richerror.Op("redispresence.GetPresence")
	// TODO - implement me
	m := make(map[uint]int64)

	for _, userID := range userIDs {
		m[userID] = timestamp.Add(time.Millisecond * -100)
	}
}
