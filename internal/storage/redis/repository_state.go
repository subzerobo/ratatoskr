package redis

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func (s redisStore) SetState(key, value string) error {
	_, err := s.redis.Set(context.Background(), key, value, 30*time.Minute).Result()
	return errors.Wrap(err, "failed to set state key")
}

func (s redisStore) GetState(key string) (string, error) {
	res, err := s.redis.Get(context.Background(), key).Result()
	return res, errors.Wrap(err, "failed to get state key")
}
