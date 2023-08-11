package dic

import (
	commonConst "pi-inventory/common/consts"
	commonService "pi-inventory/common/service"
	attachmentConst "pi-inventory/modules/attachment/consts"
	attachmentService "pi-inventory/modules/attachment/service"
	groupItemConsts "pi-inventory/modules/groupItem/consts"
	groupItemModuleController "pi-inventory/modules/groupItem/controller"
	groupItemModuleRepository "pi-inventory/modules/groupItem/repository"
	groupItemModuleService "pi-inventory/modules/groupItem/service"
	stockConsts "pi-inventory/modules/stock/consts"
	stockModuleRepository "pi-inventory/modules/stock/repository"
	stockModuleService "pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterGroupItemComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: groupItemConsts.GroupItemRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleRepository.NewGroupItemRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: groupItemConsts.GroupItemService,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleService.NewGroupItemService(
				ctn.Get(groupItemConsts.GroupItemRepository).(groupItemModuleRepository.GroupItemRepositoryInterface),
				ctn.Get(stockConsts.UnitService).(stockModuleService.UnitServiceInterface),
				ctn.Get(groupItemConsts.VariantService).(groupItemModuleService.VariantServiceInterface),
				ctn.Get(stockConsts.StockService).(stockModuleService.StockServiceInterface),
				ctn.Get(attachmentConst.AttachmentService).(attachmentService.AttachmentServiceInterface),
				ctn.Get(stockConsts.StockRepository).(stockModuleRepository.StockRepositoryInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: groupItemConsts.GroupItemController,
		Build: func(ctn di.Container) (interface{}, error) {
			return groupItemModuleController.NewGroupItemController(
				ctn.Get(groupItemConsts.GroupItemService).(groupItemModuleService.GroupItemServiceInterface),
				ctn.Get(commonConst.ActivityLogService).(commonService.ActivityLogServiceInterface),
			), nil
		},
	})
}
