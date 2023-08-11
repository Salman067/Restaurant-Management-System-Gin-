package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	compositeConst "pi-inventory/modules/composite/consts"
	compositeController "pi-inventory/modules/composite/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupCompositeRoute(_ *di.Builder, api *gin.RouterGroup) {
	CompositeController := dic.Container.Get(compositeConst.CompositeController).(compositeController.CompositeControllerInterface)
	composite := api.Group("/composite/:account_slug")
	composite.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		composite.GET("/ping", CompositeController.Pong)
		composite.POST("/create", CompositeController.CreateComposite)
		composite.GET("/list", CompositeController.FindAll)
		composite.GET("/view/:id", CompositeController.FindByID)
		composite.PUT("/update/:id", CompositeController.UpdateComposite)
	}
}
