package db

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDb() *gorm.DB {
	dbUrl := viper.GetString("DB_URL")
	opts := &gorm.Config{}

	if viper.GetString("GIN_MODE") == "debug" {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
		opts.Logger = newLogger
	}

	db, err := gorm.Open(postgres.Open(dbUrl), opts)
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(viper.GetInt("DB_MAX_CONNECTIONS"))
	sqlDb.SetMaxOpenConns(viper.GetInt("DB_MAX_CONNECTIONS"))
	// sqlDb.LogMode(viper.GetBool("DB_LOG_MODE"))

	return db
}
