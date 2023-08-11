package service

import (
	"context"
	"encoding/json"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

type NotificationServiceInterface interface {
	GetNotifications(userId uint) ([]models.Notification, error)
	SetRead(userId uint, notificationId uuid.UUID) error
	SetBatchRead(userId uint, notificationIds []uuid.UUID) error
	AddNotification(userId uint, notification *models.Notification) error
	GenerateNotification(activityLog models.ActivityLog) error
	AddBatchNotification(notifications []*models.Notification) error
	SetAllRead(userId uint) error
	GetLatestNotifications(userId uint, notificationId uuid.UUID) ([]models.Notification, error)
}

type NotificationService struct {
	ctn         di.Container
	RedisClient *redis.Client
	Db          *gorm.DB
	Logger      logger.LoggerInterface
}

func NewNotificationService(db *gorm.DB, redisClient *redis.Client, c di.Container, logger logger.LoggerInterface) NotificationServiceInterface {
	return &NotificationService{c, redisClient, db, logger}
}

func (s NotificationService) GenerateNotification(activityLog models.ActivityLog) error {
	var notifications []*models.Notification
	if len(activityLog.NotificationReceiver) != 0 {
		for _, notifiedUser := range activityLog.NotificationReceiver {
			if activityLog.CreatedBy != uint(notifiedUser.UserID) {
				notification := &models.Notification{
					ID:               uuid.New(),
					CreatedAt:        time.Now(),
					Read:             false,
					UserId:           uint(notifiedUser.UserID),
					CreatedBy:        activityLog.CreatedBy,
					AccountId:        uint64(notifiedUser.AccountID),
					ModelName:        activityLog.ModelName,
					ModelId:          activityLog.ModelId,
					ActivityType:     activityLog.ActivityType,
					SubModelName:     activityLog.SubModelName,
					SubModelId:       activityLog.SubModelId,
					NotificationText: activityLog.NotificationText,
					SearchBy:         activityLog.SearchBy,
					Notes:            activityLog.Notes,
					PreviousEntry:    activityLog.PreviousEntry,
					NewEntry:         activityLog.NewEntry,
				}

				notifications = append(notifications, notification)
			}
		}
	} else {
		var users []uint
		err := s.Db.Model(models.AccountUserPermission{}).Select("user_id").
			Where("account_id = ?", activityLog.AccountId).Scan(&users).Error
		if err != nil {
			return err
		}
		for _, user := range users {
			if user != activityLog.CreatedBy {
				notification := &models.Notification{
					ID:               uuid.New(),
					CreatedAt:        time.Now(),
					Read:             false,
					CreatedBy:        activityLog.CreatedBy,
					UserId:           user,
					AccountId:        activityLog.AccountId,
					ModelName:        activityLog.ModelName,
					ModelId:          activityLog.ModelId,
					ActivityType:     activityLog.ActivityType,
					SubModelName:     activityLog.SubModelName,
					SubModelId:       activityLog.SubModelId,
					NotificationText: activityLog.NotificationText,
					SearchBy:         activityLog.SearchBy,
					Notes:            activityLog.Notes,
					PreviousEntry:    activityLog.PreviousEntry,
					NewEntry:         activityLog.NewEntry,
				}
				notifications = append(notifications, notification)
			}
		}
	}
	return s.AddBatchNotification(notifications)
}

func (s NotificationService) AddBatchNotification(notifications []*models.Notification) error {
	for _, notification := range notifications {
		err := s.AddNotification(notification.UserId, notification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s NotificationService) AddNotification(userId uint, notification *models.Notification) error {
	key := "notification_" + strconv.Itoa(int(userId))

	notificationJson, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	_, err = s.RedisClient.LPush(context.Background(), key, notificationJson).Result()
	if err != nil {
		return err
	}

	_, err = s.RedisClient.LTrim(context.Background(), key, 0, consts.DefaultNotificationLimit-1).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s NotificationService) GetNotifications(userId uint) ([]models.Notification, error) {
	key := "notification_" + strconv.Itoa(int(userId))
	notificationsJson, err := s.RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var notifications []models.Notification
	for _, serializedItem := range notificationsJson {
		var item models.Notification
		err = json.Unmarshal([]byte(serializedItem), &item)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, item)
	}

	return notifications, nil
}

func (s NotificationService) GetLatestNotifications(userId uint, notificationId uuid.UUID) ([]models.Notification, error) {
	key := "notification_" + strconv.Itoa(int(userId))
	notificationsJson, err := s.RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var notifications []models.Notification
	for _, serializedItem := range notificationsJson {
		var notification models.Notification
		err = json.Unmarshal([]byte(serializedItem), &notification)
		if err != nil {
			return nil, err
		}
		if notification.ID != notificationId {
			notifications = append(notifications, notification)
		} else {
			break
		}
	}
	return notifications, nil
}

func (s NotificationService) SetRead(userId uint, notificationId uuid.UUID) error {
	key := "notification_" + strconv.Itoa(int(userId))
	notificationsJson, err := s.RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return err
	}

	for i, serializedItem := range notificationsJson {
		var notification models.Notification
		err = json.Unmarshal([]byte(serializedItem), &notification)
		if err != nil {
			return err
		}
		if notification.ID == notificationId {
			notification.Read = true

			newSerialized, err := json.Marshal(notification)
			if err != nil {
				return err
			}
			_, err = s.RedisClient.LSet(context.Background(), key, int64(i), newSerialized).Result()
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (s NotificationService) SetBatchRead(userId uint, notificationIds []uuid.UUID) error {
	key := "notification_" + strconv.Itoa(int(userId))
	notificationsJson, err := s.RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return err
	}

	for i, serializedItem := range notificationsJson {
		var notification models.Notification
		err = json.Unmarshal([]byte(serializedItem), &notification)
		if err != nil {
			return err
		}
		for _, notificationId := range notificationIds {
			if notification.ID == notificationId {
				notification.Read = true

				newSerialized, err := json.Marshal(notification)
				if err != nil {
					return err
				}
				_, err = s.RedisClient.LSet(context.Background(), key, int64(i), newSerialized).Result()
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func (s NotificationService) SetAllRead(userId uint) error {
	key := "notification_" + strconv.Itoa(int(userId))
	notificationsJson, err := s.RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return err
	}

	for i, serializedItem := range notificationsJson {
		var notification models.Notification
		err = json.Unmarshal([]byte(serializedItem), &notification)
		if err != nil {
			return err
		}

		if !notification.Read {
			notification.Read = true

			newSerialized, err := json.Marshal(notification)
			if err != nil {
				return err
			}
			_, err = s.RedisClient.LSet(context.Background(), key, int64(i), newSerialized).Result()
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
