package dic

import (
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
	commonConst "pi-inventory/common/consts"
	"pi-inventory/modules/attachment/consts"
	attachmentController "pi-inventory/modules/attachment/controller"
	"pi-inventory/modules/attachment/file_uploader"
	attachmentRepository "pi-inventory/modules/attachment/repository"
	attachmentService "pi-inventory/modules/attachment/service"
)

func RegisterAttachmentComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.FileUploaderFactory,
		Build: func(ctn di.Container) (interface{}, error) {
			return file_uploader.NewFileUploaderFactory(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.AttachmentRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return attachmentRepository.NewAttachmentRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.AttachmentService,
		Build: func(ctn di.Container) (interface{}, error) {
			return attachmentService.NewAttachmentService(ctn.Get(consts.AttachmentRepository).(attachmentRepository.AttachmentRepositoryInterface), ctn.Get(consts.FileUploaderFactory).(file_uploader.FileUploaderInterface)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.AttachmentController,
		Build: func(ctn di.Container) (interface{}, error) {
			return attachmentController.NewAttachmentController(ctn.Get(consts.AttachmentService).(attachmentService.AttachmentServiceInterface)), nil
		},
	})
}
