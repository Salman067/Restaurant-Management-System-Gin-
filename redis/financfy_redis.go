package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"pi-inventory/common/logger"
)

func NewRedisV8Db() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("FINANCFY_REDIS_HOST") + ":" + viper.GetString("FINANCFY_REDIS_PORT"),
		Password: viper.GetString("FINANCFY_REDIS_PASSWORD"), // no password set
		DB:       viper.GetInt("FINANCFY_REDIS_DB"),          // use default DB
	})
	if s, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.LogError("Error in petty-cash redis, ", err)
	} else {
		logger.LogInfo(s)
	}

	return redisClient
}
