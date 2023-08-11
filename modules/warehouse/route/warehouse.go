package route

import (
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	warehouseConst "pi-inventory/modules/warehouse/consts"
	warehouseController "pi-inventory/modules/warehouse/controller"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

func SetupWarehouseRoute(_ *di.Builder, api *gin.RouterGroup) {
	WarehouseController := dic.Container.Get(warehouseConst.WarehouseController).(warehouseController.WarehouseControllerInterface)
	Warehouse := api.Group("/warehouse/:account_slug")
	Warehouse.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		Warehouse.GET("/ping", WarehouseController.Pong)
		Warehouse.GET("/list", WarehouseController.FindAll)
		Warehouse.GET("/view/:id", WarehouseController.FindByID)
		Warehouse.POST("/create", WarehouseController.CreateWarehouse)
		Warehouse.PUT("/update/:id", WarehouseController.UpdateWarehouse)
		Warehouse.DELETE("/delete/:id", WarehouseController.DeleteWarehouse)
	}
}
