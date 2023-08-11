package dic

import (
	"github.com/go-redis/redis/v8"
	commonConst "pi-inventory/common/consts"
	profileCache "pi-inventory/modules/profile/cache"
	"pi-inventory/modules/profile/consts"
	profileController "pi-inventory/modules/profile/controller"
	profileRepository "pi-inventory/modules/profile/repository"
	profileService "pi-inventory/modules/profile/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterProfileComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.ProfileRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return profileRepository.NewProfileRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.ProfileCacheRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return profileCache.NewProfileCacheRepository(ctn.Get(commonConst.RedisV8DB).(*redis.Client)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.ProfileService,
		Build: func(ctn di.Container) (interface{}, error) {
			return profileService.NewProfileService(ctn.Get(consts.ProfileRepository).(profileRepository.ProfileRepositoryInterface),
				ctn.Get(consts.ProfileCacheRepository).(profileCache.ProfileCacheRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.ProfileController,
		Build: func(ctn di.Container) (interface{}, error) {
			return profileController.NewProfileController(ctn.Get(consts.ProfileService).(profileService.ProfileServiceInterface)), nil
		},
	})
}
