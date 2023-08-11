package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type StockCacheRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetAll(ctx context.Context, key string) (map[string]string, error)
}

type stockCacheRepository struct {
	RedisClient *redis.Client
}

func NewStockCacheRepository(redisClient *redis.Client) *stockCacheRepository {
	return &stockCacheRepository{redisClient}
}

func (scr *stockCacheRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := scr.RedisClient.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (scr *stockCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := scr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}
