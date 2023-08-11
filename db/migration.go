package db

import (
	attachmentSchema "pi-inventory/modules/attachment/schema"
	compositeSchema "pi-inventory/modules/composite/schema"
	groupItemSchema "pi-inventory/modules/groupItem/schema"
	"pi-inventory/modules/stock/models"
	stockSchema "pi-inventory/modules/stock/schema"
	warehouseSchema "pi-inventory/modules/warehouse/schema"

	commonModel "pi-inventory/common/models"

	"gorm.io/gorm"
)

var DemoTaxList []*models.TaxResponseBody
var DemoSupplierList []*models.SupplierResponseBody

func Migration(db *gorm.DB) {

	db.AutoMigrate(&attachmentSchema.Attachment{})
	db.AutoMigrate(&stockSchema.Category{})
	db.AutoMigrate(&stockSchema.Unit{})
	db.AutoMigrate(&warehouseSchema.Warehouse{})
	db.AutoMigrate(&groupItemSchema.Variant{})
	db.AutoMigrate(&groupItemSchema.GroupItem{})
	db.AutoMigrate(&stockSchema.Stock{})
	db.AutoMigrate(&compositeSchema.LineItem{})
	db.AutoMigrate(&compositeSchema.Composite{})
	db.AutoMigrate(&stockSchema.StockActivity{})
	db.AutoMigrate(&stockSchema.Purpose{})
	db.AutoMigrate(&commonModel.AccountUserPermission{})
}
