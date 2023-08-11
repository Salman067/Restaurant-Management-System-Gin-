package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	"pi-inventory/modules/stock/consts"
	purposeController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di/v2"
)

func SetupPurposeRoute(_ *di.Builder, api *gin.RouterGroup) {
	purposeControllers := dic.Container.Get(consts.PurposeController).(purposeController.PurposeControllerInterface)
	purpose := api.Group("/purpose/:account_slug")
	purpose.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		purpose.GET("/ping", purposeControllers.Pong)
		purpose.GET("/list", purposeControllers.FindAll)
		purpose.GET("/view/:id", purposeControllers.FindByID)
		purpose.DELETE("delete/:id", purposeControllers.DeletePurpose)
		purpose.POST("/create", purposeControllers.CreatePurpose)
		purpose.PUT("update/:id", purposeControllers.UpdatePurpose)

	}
}
