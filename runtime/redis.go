package runtime

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisScheduleStore struct {
	Client             *redis.Client
	Expiration         time.Duration
	FinishedExpiration time.Duration
}

func (rss *RedisScheduleStore) Set(s *Schedule) error {
	val, err := json.Marshal(s)
	if err != nil {
		return err
	}

	var expiration time.Duration
	if s.Finished {
		expiration = rss.FinishedExpiration
	} else {
		expiration = rss.Expiration
	}

	return rss.Client.Set(ctx, s.TraceID, val, expiration).Err()
}

func (rss *RedisScheduleStore) Get(traceID string) (*Schedule, error) {
	val, err := rss.Client.Get(ctx, traceID).Bytes()
	if err != nil {
		return nil, err
	}
	var schedule Schedule
	if err := json.Unmarshal(val, &schedule); err != nil {
		return nil, err
	}
	return &schedule, nil
}
