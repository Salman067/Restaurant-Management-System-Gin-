package main

import (
	"fmt"
	"os"
	"pi-inventory/command"
	"time"
	_ "time/tzdata"

	"gorm.io/gorm"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"

	commonConst "pi-inventory/common/consts"
	"pi-inventory/common/logger"
	"pi-inventory/config"
	DB "pi-inventory/db"
	"pi-inventory/dic"
)

func readConfig() {
	var err error
	viper.AddConfigPath(".")
	viper.SetConfigName("base")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("WARNING: file .env not found")
	} else {
		viper.SetConfigFile(".env")
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Override config parameters from environment variables if specified
	err = viper.Unmarshal(&config.Config)
	if err != nil {
		panic("error unmarshaling config")
	}

	for _, key := range viper.AllKeys() {
		viper.BindEnv(key)
	}
}

func checkDb() {
	db := dic.Container.Get(commonConst.DbService).(*gorm.DB)
	postgresDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("error connecting to db: %w", err))
	}
	err = postgresDB.Ping()
	if err != nil {
		panic(fmt.Errorf("error connecting to db: %w", err))
	}
	fmt.Println("database connected success")

	DB.Migration(db)
}

func setUTCTimezone() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		logger.LogInfo(err)
		return
	}
	time.Local = loc
}

func checkRedis() {
	redisConnectionPool := dic.Container.Get(commonConst.RedisService).(*redis.Pool)
	_, err := redisConnectionPool.Get().Do("PING")
	if err != nil {
		panic(fmt.Errorf("error connecting to redis: %w", err))
	}
}

func main() {
	setUTCTimezone()
	readConfig()
	logger.NewLogger(nil)
	dic.InitContainer()
	checkDb()
	checkRedis()
	command.Execute()
}
