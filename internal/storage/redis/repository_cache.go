package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/subzerobo/ratatoskr/internal/services/applications"
)

const (
	ApplicationCachedKey = "App:%s:AndroidParams"
)

func (s *redisStore) GetApplicationData(appUUID string) (*applications.ApplicationCachedDataModel, error) {
	key := fmt.Sprintf(ApplicationCachedKey, appUUID)
	var res applications.ApplicationCachedDataModel
	jsonData , err := s.redis.Get(context.Background(),key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonData), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *redisStore) SetApplicationData(appUUID string, model applications.ApplicationCachedDataModel) error {
	key := fmt.Sprintf(ApplicationCachedKey, appUUID)
	jsonData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = s.redis.Set(context.Background(), key, string(jsonData),0).Result()
	return err
}
