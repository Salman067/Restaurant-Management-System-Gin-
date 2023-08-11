package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	taxController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupTaxRoute(_ *di.Builder, api *gin.RouterGroup) {
	taxController := dic.Container.Get(consts.TaxController).(taxController.TaxControllerInterface)
	tax := api.Group("/tax/:account_slug")
	tax.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		tax.GET("/ping", taxController.Pong)
		tax.GET("/list", taxController.FindAll)
		tax.GET("/view/:id", taxController.FindByID)

	}
}
