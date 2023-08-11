package dic

import (
	Cache "pi-inventory/common/cache"
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	"pi-inventory/modules/stock/consts"
	stockModuleController "pi-inventory/modules/stock/controller"
	stockModuleRepository "pi-inventory/modules/stock/repository"
	stockModuleService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterCategoryComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.CategoryRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleRepository.NewCategoryRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.CategoryService,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleService.NewCategoryService(
				ctn.Get(consts.CategoryRepository).(stockModuleRepository.CategoryRepositoryInterface),
				ctn.Get(consts.StockRepository).(stockModuleRepository.StockRepositoryInterface),
				ctn.Get(commonConst.Cache).(Cache.CacheInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.CategoryController,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockModuleController.NewCategoryController(
				ctn.Get(consts.CategoryService).(stockModuleService.CategoryServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
