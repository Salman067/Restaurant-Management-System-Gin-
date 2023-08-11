package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	stockController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupStockRoute(_ *di.Builder, api *gin.RouterGroup) {
	stockControllers := dic.Container.Get(consts.StockController).(*stockController.StockController)
	stock := api.Group("/stock/:account_slug")

	stock.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		stock.GET("/ping", stockControllers.Pong)
		stock.GET("/list", stockControllers.FindAll)
		stock.GET("/view/:id", stockControllers.FindByID)
		stock.POST("/create", stockControllers.CreateStock)
		stock.PUT("/update/:id", stockControllers.UpdateStock)
		stock.DELETE("/delete/:id", stockControllers.DeleteStock)
		stock.GET("/stock-summary", stockControllers.StockSummary)
	}
}
