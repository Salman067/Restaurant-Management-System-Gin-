package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	stockConst "pi-inventory/modules/stock/consts"
	stockController "pi-inventory/modules/stock/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupCategoryRoute(_ *di.Builder, api *gin.RouterGroup) {
	categoryController := dic.Container.Get(stockConst.CategoryController).(stockController.CategoryControllerInterface)
	category := api.Group("/category/:account_slug")
	category.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		category.GET("/ping", categoryController.Pong)
		category.GET("/list", categoryController.FindAll)
		category.GET("/view/:id", categoryController.FindByID)
		category.POST("/create", categoryController.CreateCategory)
		category.PUT("/update/:id", categoryController.UpdateCategory)
		category.DELETE("/delete/:id", categoryController.DeleteCategory)
	}
}
