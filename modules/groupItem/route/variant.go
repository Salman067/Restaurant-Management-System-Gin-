package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	groupItemConst "pi-inventory/modules/groupItem/consts"
	groupItemController "pi-inventory/modules/groupItem/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupVariantRoute(_ *di.Builder, api *gin.RouterGroup) {
	variantController := dic.Container.Get(groupItemConst.VariantController).(groupItemController.VariantControllerInterface)
	variant := api.Group("/variant/:account_slug")
	variant.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		variant.GET("/list", variantController.FindAll)
		variant.GET("/view/:id", variantController.FindByID)
		variant.DELETE("delete/:id", variantController.DeleteVariant)
		variant.POST("/create", variantController.CreateVariant)
		variant.PUT("update/:id", variantController.UpdateVariant)
	}
}
