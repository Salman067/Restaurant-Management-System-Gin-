package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"pi-inventory/common/logger"
)

type SupplierCacheRepositoryInterface interface {
	GetKeys(ctx context.Context, keyPattern string) ([]string, error)
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetByFields(ctx context.Context, key string, fields []string) ([]string, error)
	GetSingleData(ctx context.Context, key string, field string) (string, error)
}

type supplierCacheRepository struct {
	RedisClient *redis.Client
}

func NewSupplierCacheRepository(redisClient *redis.Client) *supplierCacheRepository {
	return &supplierCacheRepository{redisClient}
}

func (scr *supplierCacheRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := scr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (scr *supplierCacheRepository) GetKeys(ctx context.Context, keyPattern string) ([]string, error) {
	keys, err := scr.RedisClient.Keys(ctx, keyPattern).Result()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return keys, nil
}

func (r supplierCacheRepository) GetByFields(ctx context.Context, key string, fields []string) ([]string, error) {
	logger.LogInfo("get from redis key ", key)
	pipe := r.RedisClient.Pipeline()
	cmds := map[string]*redis.StringCmd{}
	for _, field := range fields {
		cmds[field] = pipe.HGet(ctx, key, field)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	data := []string{}
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	return data, nil
}

func (r supplierCacheRepository) GetSingleData(ctx context.Context, key string, field string) (string, error) {
	logger.LogInfo("get from redis repository")
	data, err := r.RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
