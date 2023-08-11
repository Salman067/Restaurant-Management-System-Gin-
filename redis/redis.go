package redis

import (
	"pi-inventory/common/logger"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var redisConnection *redis.Pool

// var server = "invoice_redis:6379"
// var password = "invoice_1234"

func GetRedisConnection() *redis.Pool {
	password := viper.GetString("REDIS_PASS")
	server := viper.GetString("REDIS_URL")
	logger.LogInfo("found server ", server, " found password ", password)
	if redisConnection == nil {
		newRedisConnection(server, password)
	}

	s, err := redisConnection.Get().Do("PING")
	if err != nil {
		logger.LogError("Error in pi-inventory redis, ", err)
	} else {
		logger.LogInfo(s)
	}

	return redisConnection
}

func newRedisConnection(server string, password string) {
	redisConnection = &redis.Pool{
		MaxActive: viper.GetInt("REDIS_MAX_ACTIVE_CONN"),
		MaxIdle:   viper.GetInt("REDIS_MAX_IDLE_CONN"),
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			server := server
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
