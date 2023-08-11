package dic

import (
	"pi-inventory/common/consts"
	"pi-inventory/common/controller"
	commonModuleController "pi-inventory/common/controller"
	"pi-inventory/common/logger"
	"pi-inventory/common/repository"
	commonRepository "pi-inventory/common/repository"
	"pi-inventory/common/service"
	commonModuleService "pi-inventory/common/service"
	compositeConst "pi-inventory/modules/composite/consts"
	compositeService "pi-inventory/modules/composite/service"
	groupItemConst "pi-inventory/modules/groupItem/consts"
	groupItemService "pi-inventory/modules/groupItem/service"

	redisClient "github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterCommonComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.CommonService,
		Build: func(ctn di.Container) (interface{}, error) {
			return commonModuleService.NewCommonService(
				ctn.Get(compositeConst.CompositeService).(compositeService.CompositeServiceInterface),
				ctn.Get(groupItemConst.GroupItemService).(groupItemService.GroupItemServiceInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.CommonController,
		Build: func(ctn di.Container) (interface{}, error) {
			return commonModuleController.NewCommonController(ctn.Get(consts.CommonService).(commonModuleService.CommonServiceInterface)), nil
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

	// _ = builder.Add(di.Def{
	// 	Name: consts.RedisCommonService,
	// 	Build: func(ctn di.Container) (interface{}, error) {
	// 		return commonModuleService.NewRedisService(
	// 			ctn.Get(consts.RedisCommonService).(commonRepository.RedisRepositoryInterface),
	// 		), nil
	// 	},
	// })

	_ = builder.Add(di.Def{
		Name: consts.ActivityLogRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewActivityLogRepository(
				ctn.Get(consts.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.ActivityLogHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewActivityLogHandler(
				ctn.Get(consts.ActivityLogRepository).(repository.ActivityLogRepositoryInterface),
				ctn.Get(consts.NotificationService).(service.NotificationServiceInterface),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.ActivityLogService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewActivityLogService(
				ctn.Get(consts.ActivityLogHandler).(service.ActivityLogHandlerInterface),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.NotificationService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewNotificationService(ctn.Get(consts.DbService).(*gorm.DB),
				ctn.Get(consts.RedisV8DB).(*redisClient.Client),
				ctn,
				ctn.Get(consts.LoggerService).(logger.LoggerInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.NotificationController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewNotificationController(
				ctn.Get(consts.NotificationService).(service.NotificationServiceInterface),
				ctn.Get(consts.LoggerService).(logger.LoggerInterface)), nil
		},
	})
}
