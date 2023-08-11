package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type PurposeCacheRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetAll(ctx context.Context, key string) (map[string]string, error)
}

type purposeCacheRepository struct {
	RedisClient *redis.Client
}

func NewPurposeCacheRepository(redisClient *redis.Client) *purposeCacheRepository {
	return &purposeCacheRepository{redisClient}
}

func (pcr *purposeCacheRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := pcr.RedisClient.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (pcr *purposeCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := pcr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}
