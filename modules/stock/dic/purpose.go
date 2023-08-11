package dic

import (
	"github.com/go-redis/redis/v8"
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	purposeCache "pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/consts"
	purposeController "pi-inventory/modules/stock/controller"
	purposeRepository "pi-inventory/modules/stock/repository"
	purposeService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterPurposeComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.PurposeRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return purposeRepository.NewPurposeRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.PurposeCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return purposeCache.NewPurposeCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.PurposeService,
		Build: func(ctn di.Container) (interface{}, error) {
			return purposeService.NewPurposeService(ctn.Get(consts.PurposeRepository).(purposeRepository.PurposeRepositoryInterface),
				ctn.Get(consts.StockActivityRepository).(purposeRepository.StockActivityRepositoryInterface),
				ctn.Get(consts.PurposeCacheRepository).(purposeCache.PurposeCacheRepositoryInterface)), nil

		},
	})

	_ = builder.Add(di.Def{
		Name: consts.PurposeController,
		Build: func(ctn di.Container) (interface{}, error) {
			return purposeController.NewPurposeController(
				ctn.Get(consts.PurposeService).(purposeService.PurposeServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
