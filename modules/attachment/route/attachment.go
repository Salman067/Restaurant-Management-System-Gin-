package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/attachment/consts"
	attachmentController "pi-inventory/modules/attachment/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupAttachmentRoute(_ *di.Builder, api *gin.RouterGroup) {
	fileController := dic.Container.Get(consts.AttachmentController).(attachmentController.AttachmentControllerInterface)
	file := api.Group("/attachment/:account_slug")
	file.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		file.GET("", fileController.GetSingleFile)
		file.DELETE("", fileController.DeleteSingleFile)
		file.POST("/upload/single", fileController.UploadSingleFile)
		file.GET("/single/:path", fileController.GetSingleAttachmentFile)
		file.GET("/list/:attachmentKey", fileController.FetchAttachments)
		file.POST("/store-attachment-paths", fileController.UploadMultipleAttachment)
	}
}
