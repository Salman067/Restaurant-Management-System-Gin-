package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RedisUserInfo struct {
	ID            uint
	Name          string         `json:"name"`
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	ID               uuid.UUID
	CreatedAt        time.Time `json:"created_at"`
	Read             bool      `json:"read"`
	UserId           uint      `json:"user_id"`           // the user that will receive the notification
	CreatedBy        uint      `json:"created_by"`        // the user that created the notification
	AccountId        uint64    `json:"account_id"`        // related account for this notification
	PreviousEntry    string    `json:"previous_entry"`    // previous models value before action
	NewEntry         string    `json:"new_entry"`         // post models value after action
	ModelName        string    `json:"model_name"`        // models name like category,paymentmode
	ModelId          uint      `json:"model_id"`          // id of models
	ActivityType     string    `json:"activity_type"`     // models create,update,delete
	SubModelName     string    `json:"sub_model_name"`    // models name like category,paymentmode
	SubModelId       uint      `json:"sub_model_id"`      // models id
	NotificationText string    `json:"notification_text"` // text will be shown in the notification
	SearchBy         string    `json:"search_by"`         // can search by to find
	Notes            string    `json:"notes"`             // additional notes
	Link             string    `json:"link"`
}

type ActivityLog struct {
	gorm.Model
	CreatedBy            uint                    `json:"created_by"` // action performer user_id
	AccountId            uint64                  `json:"account_id"`
	PreviousEntry        string                  `json:"previous_entry"`        // previous models value before action
	NewEntry             string                  `json:"new_entry"`             // post models value after action
	ModelName            string                  `json:"model_name"`            //models name like category,paymentmode
	ModelId              uint                    `json:"model_id"`              // id of models
	ActivityType         string                  `json:"activity_type"`         //models create,update,delete
	SubModelName         string                  `json:"sub_model_name"`        //models name like category,paymentmode
	SubModelId           uint                    `json:"sub_model_id"`          // models id
	NotificationText     string                  `json:"notification_text"`     // text will be shown in the notification
	SearchBy             string                  `json:"search_by"`             // can search by to find
	Notes                string                  `json:"notes"`                 // additional notes
	PushNotification     bool                    `json:"push_notification"`     // will notification Will be pushed
	NotificationReceiver []AccountUserPermission `json:"notification_receiver"` // will get the notification. empty means all user of that account
	// PreviousModel    interface{} `gorm:"->;-:migration"`
	// NewModel         interface{} `gorm:"->;-:migration"`
}
