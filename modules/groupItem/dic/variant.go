package dic

import (
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	groupItemConsts "pi-inventory/modules/groupItem/consts"
	groupItemModuleController "pi-inventory/modules/groupItem/controller"
	groupItemModuleRepository "pi-inventory/modules/groupItem/repository"
	groupItemModuleService "pi-inventory/modules/groupItem/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterVariantComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: groupItemConsts.VariantRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleRepository.NewVariantRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: groupItemConsts.VariantService,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleService.NewVariantService(ctn.Get(groupItemConsts.VariantRepository).(groupItemModuleRepository.VariantRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: groupItemConsts.VariantController,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleController.NewVariantController(
				ctn.Get(groupItemConsts.VariantService).(groupItemModuleService.VariantServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
