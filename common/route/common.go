package route

import (
	"pi-inventory/common/consts"
	"pi-inventory/common/controller"
	"pi-inventory/dic"
	"pi-inventory/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupCommonRoute(_ *di.Builder, api *gin.RouterGroup) {
	commonController := dic.Container.Get(consts.CommonController).(controller.CommonControllerInterface)
	common := api.Group("/group-and-composite-item/:account_slug")
	common.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{

		common.GET("/ping", commonController.Pong)
		common.GET("/summary", commonController.GroupItemAndCompositeItemSummary)
	}
}
