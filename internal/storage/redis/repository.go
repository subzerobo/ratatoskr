package redis

import "github.com/go-redis/redis/v8"

type redisStore struct {
	redis *redis.Client
}

func CreateRedisStore(redis *redis.Client) *redisStore {
	store := &redisStore{
		redis: redis,
	}
	return store
}
