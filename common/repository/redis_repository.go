package repository

import (
	"context"
	"pi-inventory/common/logger"

	"github.com/go-redis/redis/v8"
)

type RedisRepositoryInterface interface {
	Set(ctx context.Context, key string, value map[string]interface{}) error
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetByKey(ctx context.Context, key string) (string, error)
	GetSingleData(ctx context.Context, key string, field string) (string, error)
	GetByFields(ctx context.Context, key string, fields []string) ([]string, error)
	Delete(ctx context.Context, key string, fields []string) error
}

type RedisRepository struct {
	RedisClient *redis.Client
	logger      logger.LoggerInterface
}

func NewRedisRepository(redisClient *redis.Client, logger logger.LoggerInterface) RedisRepositoryInterface {
	return &RedisRepository{redisClient, logger}
}

func (r RedisRepository) Set(ctx context.Context, key string, value map[string]interface{}) error {
	logger.LogInfo("set to redis ", key)
	err := r.RedisClient.HMSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) Delete(ctx context.Context, key string, fields []string) error {
	var err error
	if len(fields) == 0 || fields == nil {
		err = r.RedisClient.Del(ctx, key).Err()
	} else {
		err = r.RedisClient.HDel(ctx, key, fields...).Err()
	}
	if err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) GetAll(ctx context.Context, key string) (map[string]string, error) {
	logger.LogInfo("key ", key)
	data, err := r.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RedisRepository) GetSingleData(ctx context.Context, key string, field string) (string, error) {
	logger.LogInfo("get tax from redis repository ", key, " ", field)
	data, err := r.RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r RedisRepository) GetByFields(ctx context.Context, key string, fields []string) ([]string, error) {
	logger.LogInfo("key ", key)
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

func (r RedisRepository) GetByKey(ctx context.Context, key string) (string, error) {
	logger.LogInfo("key ", key)
	data, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
