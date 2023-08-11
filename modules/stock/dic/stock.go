package dic

import (
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	attachmentConst "pi-inventory/modules/attachment/consts"
	attachmentService "pi-inventory/modules/attachment/service"
	stockCache "pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/consts"
	stockController "pi-inventory/modules/stock/controller"
	stockRepository "pi-inventory/modules/stock/repository"
	stockService "pi-inventory/modules/stock/service"

	"github.com/go-redis/redis/v8"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterStockComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.StockRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockRepository.NewStockRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.StockCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockCache.NewStockCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.StockService,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockService.NewStockService(ctn.Get(consts.StockRepository).(stockRepository.StockRepositoryInterface),
				ctn.Get(attachmentConst.AttachmentService).(attachmentService.AttachmentServiceInterface),
				ctn.Get(consts.CategoryService).(stockService.CategoryServiceInterface),
				ctn.Get(consts.UnitService).(stockService.UnitServiceInterface),
				ctn.Get(consts.StockCacheRepository).(stockCache.StockCacheRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.StockController,
		Build: func(ctn di.Container) (interface{}, error) {
			return stockController.NewStockController(ctn.Get(consts.StockService).(stockService.StockServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
