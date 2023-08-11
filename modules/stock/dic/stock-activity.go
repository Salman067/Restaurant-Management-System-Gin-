package dic

import (
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	stockConst "pi-inventory/modules/stock/consts"
	stockController "pi-inventory/modules/stock/controller"
	stockRepository "pi-inventory/modules/stock/repository"
	stockService "pi-inventory/modules/stock/service"
	warehouseConst "pi-inventory/modules/warehouse/consts"
	warehouseService "pi-inventory/modules/warehouse/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterStockActivityComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: stockConst.StockActivityRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockRepository.NewStockActivityRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: stockConst.StockActivityService,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockService.NewStockActivityService(ctn.Get(stockConst.StockActivityRepository).(stockRepository.StockActivityRepositoryInterface),
				ctn.Get(stockConst.StockService).(stockService.StockServiceInterface),
				ctn.Get(stockConst.PurposeService).(stockService.PurposeServiceInterface),
				ctn.Get(warehouseConst.WarehouseService).(warehouseService.WarehouseServiceInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: stockConst.StockActivityController,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockController.NewStockActivityController(ctn.Get(stockConst.StockActivityService).(stockService.StockActivityServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
