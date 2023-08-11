package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type TaxCacheRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetKeys(ctx context.Context, keyPattern string) ([]string, error)
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetSingleData(ctx context.Context, key string, field string) (string, error)
}

type taxCacheRepository struct {
	RedisClient *redis.Client
}

func NewTaxCacheRepository(redisClient *redis.Client) *taxCacheRepository {
	return &taxCacheRepository{redisClient}
}

func (tcr *taxCacheRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := tcr.RedisClient.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (tcr *taxCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := tcr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tcr *taxCacheRepository) GetKeys(ctx context.Context, keyPattern string) ([]string, error) {
	keys, err := tcr.RedisClient.Keys(ctx, keyPattern).Result()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return keys, nil
}

func (r taxCacheRepository) GetSingleData(ctx context.Context, key string, field string) (string, error) {
	logger.LogInfo("get from redis repository")
	data, err := r.RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
