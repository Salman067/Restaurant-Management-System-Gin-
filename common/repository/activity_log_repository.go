package repository

import (
	"pi-inventory/common/logger"
	"pi-inventory/common/models"

	"gorm.io/gorm"
)

type ActivityLogRepositoryInterface interface {
	Create(log models.ActivityLog) error
}

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepositoryInterface {
	return &ActivityLogRepository{
		db: db,
	}
}

func (activityLogRepo *ActivityLogRepository) Create(log models.ActivityLog) error {
	logger.LogInfo(log)
	//return activityLogRepo.db.Create(&log).Error
	return nil
}
