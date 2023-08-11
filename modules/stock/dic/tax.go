package dic

import (
	"github.com/go-redis/redis/v8"
	commonConst "pi-inventory/common/consts"
	taxCache "pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/consts"
	taxController "pi-inventory/modules/stock/controller"
	taxRepository "pi-inventory/modules/stock/repository"
	taxService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterTaxComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.TaxRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return taxRepository.NewTaxRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.TaxCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return taxCache.NewTaxCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.TaxService,
		Build: func(ctn di.Container) (interface{}, error) {
			return taxService.NewTaxService(ctn.Get(consts.TaxRepository).(taxRepository.TaxRepositoryInterface),
				ctn.Get(consts.TaxCacheRepository).(taxCache.TaxCacheRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.TaxController,
		Build: func(ctn di.Container) (interface{}, error) {
			return taxController.NewTaxController(ctn.Get(consts.TaxService).(taxService.TaxServiceInterface)), nil
		},
	})
}
