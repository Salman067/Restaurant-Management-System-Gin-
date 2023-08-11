package dic

import (
	commonConst "pi-inventory/common/consts"
	attachmentConst "pi-inventory/modules/attachment/consts"
	attachmentService "pi-inventory/modules/attachment/service"
	"pi-inventory/modules/composite/consts"
	compositeModuleController "pi-inventory/modules/composite/controller"
	compositeModuleRepository "pi-inventory/modules/composite/repository"
	compositeModuleService "pi-inventory/modules/composite/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterCompositeComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.CompositeRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return compositeModuleRepository.NewCompositeRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.CompositeService,
		Build: func(ctn di.Container) (interface{}, error) {
			return compositeModuleService.NewCompositeService(
				ctn.Get(consts.CompositeRepository).(compositeModuleRepository.CompositeRepositoryInterface),
				ctn.Get(consts.LineItemService).(compositeModuleService.LineItemServiceInterface),
				ctn.Get(attachmentConst.AttachmentService).(attachmentService.AttachmentServiceInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.CompositeController,
		Build: func(ctn di.Container) (interface{}, error) {
			return compositeModuleController.NewCompositeController(ctn.Get(consts.CompositeService).(compositeModuleService.CompositeServiceInterface)), nil
		},
	})
}
