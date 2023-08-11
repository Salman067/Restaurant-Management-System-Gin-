package cache

import (
	"errors"
	"fmt"
	"pi-inventory/common/logger"

	"github.com/gomodule/redigo/redis"
)

type CacheInterface interface {
	Set(key string, val string) error
	Get(key string) (string, error)
	Remove(key string) error
}

type cache struct {
	redisPool *redis.Pool
}

func NewCache(redisPool *redis.Pool) *cache {
	return &cache{redisPool: redisPool}
}

func (uc *cache) Set(key string, val string) error {
	conn := uc.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val, "EX", 1800)
	if err != nil {
		logger.LogError(fmt.Sprintf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error()))
		return err
	}
	return nil
}

func (uc *cache) Get(key string) (string, error) {
	conn := uc.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		logger.LogError(fmt.Sprintf("ERROR: fail get key %s, error %s", key, err.Error()))
		return "", err
	}

	return s, nil
}
func (uc *cache) Remove(key string) error {
	conn := uc.redisPool.Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		logger.LogError(fmt.Sprintf("ERROR: fail delete key %s, error %s", key, err.Error()))
		return err
	}
	if res == 1 {
		logger.LogInfo("Key delete successfully.")
		return nil
	} else if res == 0 {
		logger.LogError("Key not found.")
		return errors.New("KeyNotFound")
	} else {
		logger.LogError("Unexpected result:", res)
		return errors.New("UnexpectedResult")
	}
}
