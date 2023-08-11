package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	supplierController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupSupplierRoute(_ *di.Builder, api *gin.RouterGroup) {
	supplierController := dic.Container.Get(consts.SupplierController).(supplierController.SupplierControllerInterface)
	supplier := api.Group("/supplier/:account_slug")
	supplier.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		supplier.GET("/ping", supplierController.Pong)
		supplier.GET("/list", supplierController.FindAll)
		supplier.GET("/view/:id", supplierController.FindByID)
	}
}
