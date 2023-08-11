package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/groupItem/consts"
	groupItemController "pi-inventory/modules/groupItem/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupGroupItemRoute(_ *di.Builder, api *gin.RouterGroup) {
	groupItemController := dic.Container.Get(consts.GroupItemController).(groupItemController.GroupItemControllerInterface)
	groupItem := api.Group("/group-item/:account_slug")
	groupItem.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		groupItem.GET("/ping", groupItemController.Pong)
		groupItem.POST("/create", groupItemController.CreateGroupItem)
		groupItem.GET("/view/:id", groupItemController.FindByID)
		groupItem.GET("/list", groupItemController.FindAll)
		groupItem.PUT("/update/:id", groupItemController.UpdateGroupItem)
		// groupItem.DELETE("/delete/:id", groupItemController.DeleteGroupItem)
	}
}
