package dic

import (
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	stockConst "pi-inventory/modules/stock/consts"
	stockModuleRepository "pi-inventory/modules/stock/repository"
	warehouseModuleCacheRepository "pi-inventory/modules/warehouse/cache"
	"pi-inventory/modules/warehouse/consts"
	warehouseModuleController "pi-inventory/modules/warehouse/controller"
	warehouseModuleRepository "pi-inventory/modules/warehouse/repository"
	warehouseModuleService "pi-inventory/modules/warehouse/service"
)

func RegisterWarehouseComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.WarehouseRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return warehouseModuleRepository.NewWarehouseRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.WarehouseCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return warehouseModuleCacheRepository.NewWarehouseCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.WarehouseService,
		Build: func(ctn di.Container) (interface{}, error) {
			return warehouseModuleService.NewWarehouseService(
				ctn.Get(consts.WarehouseRepository).(warehouseModuleRepository.WarehouseRepositoryInterface),
				ctn.Get(stockConst.StockRepository).(stockModuleRepository.StockRepositoryInterface),
				ctn.Get(stockConst.StockActivityRepository).(stockModuleRepository.StockActivityRepositoryInterface),
				ctn.Get(consts.WarehouseCacheRepository).(warehouseModuleCacheRepository.WarehouseCacheRepositoryInterface),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.WarehouseController,
		Build: func(ctn di.Container) (interface{}, error) {
			return warehouseModuleController.NewWarehouseController(
				ctn.Get(consts.WarehouseService).(warehouseModuleService.WarehouseServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
