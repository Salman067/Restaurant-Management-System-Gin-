package dic

import (
	"github.com/go-redis/redis/v8"
	commonConst "pi-inventory/common/consts"
	supplierCache "pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/consts"
	supplierController "pi-inventory/modules/stock/controller"
	supplierRepository "pi-inventory/modules/stock/repository"
	supplierService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterSupplierComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.SupplierRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return supplierRepository.NewSupplierRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.SupplierCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return supplierCache.NewSupplierCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.SupplierService,
		Build: func(ctn di.Container) (interface{}, error) {
			return supplierService.NewSupplierService(ctn.Get(consts.SupplierRepository).(supplierRepository.SupplierRepositoryInterface),
				ctn.Get(consts.SupplierCacheRepository).(supplierCache.SupplierCacheRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.SupplierController,
		Build: func(ctn di.Container) (interface{}, error) {
			return supplierController.NewSupplierController(ctn.Get(consts.SupplierService).(supplierService.SupplierServiceInterface)), nil
		},
	})
}
