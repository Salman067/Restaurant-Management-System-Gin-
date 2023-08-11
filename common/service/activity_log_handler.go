package service

import (
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	"pi-inventory/common/repository"
)

type ActivityLogHandlerInterface interface {
	Create(log models.ActivityLog) error
	HandleBatchActivityLog(logs []models.ActivityLog) error
}

type activityLogHandler struct {
	activityLogRepository repository.ActivityLogRepositoryInterface
	notificationService   NotificationServiceInterface
}

func NewActivityLogHandler(activityLogRepository repository.ActivityLogRepositoryInterface, notificationService NotificationServiceInterface) ActivityLogHandlerInterface {
	return &activityLogHandler{
		activityLogRepository: activityLogRepository,
		notificationService:   notificationService,
	}
}

// gorm.Model
//
//	ModelId          uint   `json:"model_id"`          // id of models
//	ActivityType     string `json:"activity_type"`     //models create,update,delete
//	SubModelName     string `json:"sub_model_name"`    //models name like category,paymentmode
//	SubModelId       uint   `json:"sub_model_id"`      // models id
//	NotificationText string `json:"notification_text"` // text will be shown in the notification
//	SearchBy         string `json:"search_by"`         // can search by to find
//	Notes            string `json:"notes"`             // additional notes
//	PushNotification bool   `json:"push_notification"` // will notification Will be pushed
//
// newModel interface{},
//
//	previousModel interface{},
//	userId uint,
//	modelName string,
//	activityType string

func (handler activityLogHandler) Create(log models.ActivityLog) error {

	// newModel := log.NewModel
	// previousModel := log.PreviousModel
	// modelName := log.ModelName
	// activityType := log.ActivityType

	// if activityType == consts.ActionCreate || activityType == consts.ActionUpdate {
	// newModelJson, _ := json.Marshal(newModel)
	// log.NewEntry = string(newModelJson)
	// log.NewEntry = log.NewEntry
	// }
	// if activityType == consts.ActionDelete || activityType == consts.ActionUpdate {
	// 	previousModelJson, _ := json.Marshal(previousModel)
	// 	log.PreviousEntry = string(previousModelJson)
	// }
	// if modelName == consts.ModelNameCategory {
	// 	// var category *entity.Category
	// 	if activityType == consts.ActionDelete {
	// 		// category = previousModel.(*entity.Category)
	// 		log.NotificationText = "A category has been deleted from your cash book"
	// 	} else if activityType == consts.ActionUpdate {
	// 		// category = newModel.(*entity.Category)
	// 		log.NotificationText = "A category has been updated from your cash book"
	// 	} else {
	// 		// category = newModel.(*entity.Category)
	// 		log.NotificationText = "A category has been created from your cash book"
	// 	}
	// 	log.ModelId = 10  // category.ID
	// 	log.AccountId = 1 // uint64(category.AccountId)
	// 	log.PushNotification = true
	// }
	// else if modelName == consts.ModelNamePaymentMode {
	// 	var paymentMode *entity.PaymentMode
	// 	if activityType == consts.ActionDelete {
	// 		// paymentMode = previousModel.(*entity.PaymentMode)
	// 		log.NotificationText = "A payment mode has been deleted from your cash book"
	// 	} else if activityType == consts.ActionUpdate {
	// 		paymentMode = newModel.(*entity.PaymentMode)
	// 		log.NotificationText = "A payment mode has been updated from your cash book"
	// 	} else {
	// 		paymentMode = newModel.(*entity.PaymentMode)
	// 		log.NotificationText = "A payment mode has been created from your cash book"
	// 	}
	// 	log.ModelId = paymentMode.ID
	// 	log.AccountId = uint64(paymentMode.AccountId)
	// 	log.PushNotification = true
	// } else if modelName == consts.ModelNameContact {
	// 	var contact *entity.Contact
	// 	if activityType == consts.ActionDelete {
	// 		// contact = previousModel.(*entity.Contact)
	// 		log.NotificationText = "A contact has been deleted from your cash book"
	// 	} else if activityType == consts.ActionUpdate {
	// 		contact = newModel.(*entity.Contact)
	// 		log.NotificationText = "A contact has been updated from your cash book"
	// 	} else {
	// 		contact = newModel.(*entity.Contact)
	// 		log.NotificationText = "A contact has been created from your cash book"
	// 	}
	// 	log.ModelId = contact.ID
	// 	log.AccountId = uint64(contact.AccountId)
	// 	log.PushNotification = false
	// } else if modelName == consts.ModelNameAccount {
	// 	var account *entity.Account
	// 	if activityType == consts.ActionDelete {
	// 		account = previousModel.(*entity.Account)
	// 		log.NotificationText = "A cash book has been deleted"
	// 	} else if activityType == consts.ActionUpdate {
	// 		account = newModel.(*entity.Account)
	// 		log.NotificationText = "A cash book has been updated"
	// 	} else {
	// 		account = newModel.(*entity.Account)
	// 		log.NotificationText = "A cash book has been created"
	// 	}
	// 	log.ModelId = account.ID
	// 	log.AccountId = uint64(account.ID)
	// 	log.PushNotification = true
	// }
	// log.NewModel = nil
	// log.PreviousModel = nil
	if log.PushNotification {
		err := handler.notificationService.GenerateNotification(log)
		if err != nil {
			logger.LogError(err)
		}
	}
	return handler.activityLogRepository.Create(log)
}

// func (handler activityLogHandler) CreateActivityLogForDeleteModel(log entity.ActivityLog) error {
// 	return handler.activityLogService.CreateActivityLog(log)
// }
// func (handler activityLogHandler) CreateActivityLogForUpdateModel(log entity.ActivityLog) error {
// 	return handler.activityLogService.CreateActivityLog(log)
// }

func (handler activityLogHandler) HandleBatchActivityLog(logs []models.ActivityLog) error {
	for _, log := range logs {
		if log.PushNotification {
			err := handler.notificationService.GenerateNotification(log)
			if err != nil {
				logger.LogError(err)
			}
		}
	}
	//TODO implement me
	panic("implement me")
}
