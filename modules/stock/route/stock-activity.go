package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	stockActivityController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupStockActivityRoute(_ *di.Builder, api *gin.RouterGroup) {
	stockActivityControllers := dic.Container.Get(consts.StockActivityController).(stockActivityController.StockActivityControllerInterface)
	stockActivity := api.Group("/stock-activity/:account_slug")
	stockActivity.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		stockActivity.GET("/ping", stockActivityControllers.Pong)
		stockActivity.POST("/create", stockActivityControllers.CreateStockActivity)
		stockActivity.GET("/list", stockActivityControllers.FindAll)
		stockActivity.POST("/bulk-create", stockActivityControllers.CreateBulkStockActivity)
	}
}
