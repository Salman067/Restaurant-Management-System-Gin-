package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type WarehouseCacheRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetAll(ctx context.Context, key string) (map[string]string, error)
}

type warehouseCacheRepository struct {
	RedisClient *redis.Client
}

func NewWarehouseCacheRepository(redisClient *redis.Client) *warehouseCacheRepository {
	return &warehouseCacheRepository{redisClient}
}

func (wcr *warehouseCacheRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := wcr.RedisClient.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (wcr *warehouseCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := wcr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}
