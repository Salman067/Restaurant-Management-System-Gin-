package service

import (
	"context"
	"encoding/json"
	"fmt"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"

	"github.com/go-redis/redis/v8"
)

type RedisServiceInterface interface {
	GetRedisAccountInfo(accountInfo *models.RedisAccountInfo, accountSlug string) (*models.RedisAccountInfo, error)

	// UpdateLastVoucherNumberForAccount(accountSlug string, voucherNumber uint64) error
}

type RedisService struct {
	RedisClient *redis.Client
	Logger      logger.LoggerInterface
}

func NewRedisService(redisClient *redis.Client, logger logger.LoggerInterface) RedisServiceInterface {
	return &RedisService{redisClient, logger}
}

func (r *RedisService) GetRedisAccountInfo(accountInfo *models.RedisAccountInfo, accountSlug string) (*models.RedisAccountInfo, error) {
	accountInfoJson, err := r.RedisClient.Get(context.Background(), accountSlug).Result()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal([]byte(accountInfoJson), &accountInfo)
	if err != nil {
		logger.LogError(err)
	}

	accountInfo.UpdatePermission(accountInfo)

	return accountInfo, nil
}
