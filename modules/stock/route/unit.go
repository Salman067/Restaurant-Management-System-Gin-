package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	unitController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupUnitRoute(_ *di.Builder, api *gin.RouterGroup) {
	unitController := dic.Container.Get(consts.UnitController).(unitController.UnitControllerInterface)
	unit := api.Group("/unit/:account_slug")
	unit.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		unit.GET("/view/:id", unitController.FindByID)
		unit.GET("/list", unitController.FindAll)
		unit.DELETE("delete/:id", unitController.DeleteUnit)
		unit.POST("/create", unitController.CreateUnit)
		unit.PUT("update/:id", unitController.UpdateUnit)
	}
}
