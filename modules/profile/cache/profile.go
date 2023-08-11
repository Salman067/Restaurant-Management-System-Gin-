package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type ProfileCacheRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetKeys(ctx context.Context, keyPattern string) ([]string, error)
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetSingleData(ctx context.Context, key string, field string) (string, error)
}

type profileCacheRepository struct {
	RedisClient *redis.Client
}

func NewProfileCacheRepository(redisClient *redis.Client) *profileCacheRepository {
	return &profileCacheRepository{redisClient}
}

func (pcr *profileCacheRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := pcr.RedisClient.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (pcr *profileCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := pcr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (pcr *profileCacheRepository) GetKeys(ctx context.Context, keyPattern string) ([]string, error) {
	keys, err := pcr.RedisClient.Keys(ctx, keyPattern).Result()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return keys, nil
}

func (pcr profileCacheRepository) GetSingleData(ctx context.Context, key string, field string) (string, error) {
	logger.LogInfo("get from redis repository")
	data, err := pcr.RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
