package dic

import (
	Cache "pi-inventory/common/cache"
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	stockConsts "pi-inventory/modules/stock/consts"
	stockModuleController "pi-inventory/modules/stock/controller"
	stockModuleRepository "pi-inventory/modules/stock/repository"
	stockModuleService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterUnitComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: stockConsts.UnitRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleRepository.NewUnitRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: stockConsts.UnitService,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleService.NewUnitService(
				ctn.Get(stockConsts.UnitRepository).(stockModuleRepository.UnitRepositoryInterface),
				ctn.Get(stockConsts.StockRepository).(stockModuleRepository.StockRepositoryInterface),
				ctn.Get(commonConst.Cache).(Cache.CacheInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: stockConsts.UnitController,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleController.NewUnitController(
				ctn.Get(stockConsts.UnitService).(stockModuleService.UnitServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
