package dic

import (
	commonConst "pi-inventory/common/consts"
	"pi-inventory/modules/composite/consts"
	lineItemModuleRepository "pi-inventory/modules/composite/repository"
	lineItemModuleService "pi-inventory/modules/composite/service"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/service"

	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterlineItemComponent(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: consts.LineItemRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return lineItemModuleRepository.NewLineItemRepository(ctn.Get(commonConst.DbService).(*gorm.DB)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: consts.LineItemService,
		Build: func(ctn di.Container) (interface{}, error) {
			return lineItemModuleService.NewLineItemService(ctn.Get(consts.LineItemRepository).(lineItemModuleRepository.LineItemRepositoryInterface), ctn.Get(stockConst.StockService).(service.StockServiceInterface)), nil
		},
	})
}
