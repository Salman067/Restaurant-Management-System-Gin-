package dic

import (
	Cache "pi-inventory/common/cache"
	"pi-inventory/common/consts"
	commonContainer "pi-inventory/common/dic"
	"pi-inventory/common/logger"
	"pi-inventory/db"
	attachmentContainer "pi-inventory/modules/attachment/dic"
	compositeContainer "pi-inventory/modules/composite/dic"
	groupItemContainer "pi-inventory/modules/groupItem/dic"
	profileContainer "pi-inventory/modules/profile/dic"
	stockContainer "pi-inventory/modules/stock/dic"
	warehouseContainer "pi-inventory/modules/warehouse/dic"
	redisModule "pi-inventory/redis"

	commonRepository "pi-inventory/common/repository"
	commonModuleService "pi-inventory/common/service"

	"github.com/getsentry/raven-go"
	"github.com/gomodule/redigo/redis"
	"github.com/sarulabs/di/v2"

	redisClient "github.com/go-redis/redis/v8"
)

var CommonBuilder *di.Builder
var Container di.Container

func InitContainer() di.Container {
	builder := InitBuilder()
	Container = builder.Build()
	return Container
}

func InitBuilder() *di.Builder {
	CommonBuilder, _ = di.NewBuilder()
	RegisterCommonServices(CommonBuilder)
	attachmentContainer.RegisterAttachmentComponent(CommonBuilder)
	stockContainer.RegisterCategoryComponent(CommonBuilder)
	stockContainer.RegisterUnitComponent(CommonBuilder)
	warehouseContainer.RegisterWarehouseComponent(CommonBuilder)
	groupItemContainer.RegisterVariantComponent(CommonBuilder)
	groupItemContainer.RegisterGroupItemComponent(CommonBuilder)
	stockContainer.RegisterStockComponent(CommonBuilder)
	compositeContainer.RegisterlineItemComponent(CommonBuilder)
	compositeContainer.RegisterCompositeComponent(CommonBuilder)
	stockContainer.RegisterStockActivityComponent(CommonBuilder)
	stockContainer.RegisterPurposeComponent(CommonBuilder)
	stockContainer.RegisterTaxComponent(CommonBuilder)
	stockContainer.RegisterSupplierComponent(CommonBuilder)
	commonContainer.RegisterCommonComponent(CommonBuilder)
	profileContainer.RegisterProfileComponent(CommonBuilder)

	return CommonBuilder
}

func RegisterCommonServices(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.RavenClient,
		Build: func(ctn di.Container) (interface{}, error) {
			return logger.NewRavenClient(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.LoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return logger.NewLogger(ctn.Get(consts.RavenClient).(*raven.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.DbService,
		Build: func(ctn di.Container) (interface{}, error) {
			return db.NewGormDb(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.RedisService,
		Build: func(ctn di.Container) (interface{}, error) {
			return redisModule.GetRedisConnection(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.Cache,
		Build: func(ctn di.Container) (interface{}, error) {
			return Cache.NewCache(ctn.Get(consts.RedisService).(*redis.Pool)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.RedisV8DB,
		Build: func(ctn di.Container) (interface{}, error) {
			return redisModule.NewRedisV8Db(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.RedisCommonRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return commonRepository.NewRedisRepository(
				ctn.Get(consts.RedisV8DB).(*redisClient.Client),
				ctn.Get(consts.LoggerService).(logger.LoggerInterface),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.RedisCommonService,
		Build: func(ctn di.Container) (interface{}, error) {
			return commonModuleService.NewRedisService(
				ctn.Get(consts.RedisV8DB).(*redisClient.Client),
				ctn.Get(consts.LoggerService).(logger.LoggerInterface)), nil
		},
	})

	// _ = builder.Add(di.Def{
	// 	Name: consts.RedisCommonService,
	// 	Build: func(ctn di.Container) (interface{}, error) {
	// 		return commonModuleService.NewRedisService(
	// 			ctn.Get(consts.RedisCommonService).(commonRepository.RedisRepositoryInterface),
	// 		), nil
	// 	},
	// })

}
