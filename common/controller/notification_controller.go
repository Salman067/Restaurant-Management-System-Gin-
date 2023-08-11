package controller

import (
	"net/http"
	"pi-inventory/common/logger"
	"pi-inventory/common/service"
	"pi-inventory/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type NotificationController struct {
	*BaseController
	service service.NotificationServiceInterface
	logger  logger.LoggerInterface
}

func NewNotificationController(service service.NotificationServiceInterface, logger logger.LoggerInterface) *NotificationController {
	controller := NewBaseController(logger)
	return &NotificationController{
		BaseController: controller,
		service:        service,
		logger:         logger,
	}
}

func (c NotificationController) GetNotifications(context *gin.Context) {
	userId := context.GetInt64("user_id")

	resp, err := c.service.GetNotifications(uint(userId))
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, resp)
	return
}

func (c NotificationController) MarkRead(context *gin.Context) {
	userId := context.GetInt64("user_id")
	notificationId, err := uuid.Parse(context.Params.ByName("id"))

	err = c.service.SetRead(uint(userId), notificationId)
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, utils.Trans("notificationReadSuccessfully", nil))
	return
}

func (c NotificationController) MarkReadBatch(context *gin.Context) {
	userId := context.GetInt64("user_id")

	type BatchRead struct {
		NotificationIds []uuid.UUID `json:"notification_ids" binding:"required"`
	}
	batchRead := BatchRead{}
	if err := context.ShouldBindBodyWith(&batchRead, binding.JSON); err != nil {
		c.ReplyValidationError(context, err)
		return
	}

	err := c.service.SetBatchRead(uint(userId), batchRead.NotificationIds)
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, utils.Trans("notificationReadSuccessfully", nil))
	return
}

func (c NotificationController) MarkReadAll(context *gin.Context) {
	userId := context.GetInt64("user_id")

	err := c.service.SetAllRead(uint(userId))
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, utils.Trans("notificationReadSuccessfully", nil))
	return
}

func (c NotificationController) GetLatestNotifications(context *gin.Context) {
	userId := context.GetInt64("user_id")
	notificationId, err := uuid.Parse(context.Params.ByName("id"))
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusBadRequest)
		return
	}

	resp, err := c.service.GetLatestNotifications(uint(userId), notificationId)
	if err != nil {
		c.ReplyError(context, utils.Trans("somethingWentWrong", nil), http.StatusInternalServerError)
		return
	}

	c.ReplySuccess(context, resp)
	return
}
