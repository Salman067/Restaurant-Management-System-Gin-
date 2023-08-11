package db

import (
	attachmentSchema "pi-inventory/modules/attachment/schema"
	compositeSchema "pi-inventory/modules/composite/schema"
	groupItemSchema "pi-inventory/modules/groupItem/schema"
	"pi-inventory/modules/stock/models"
	stockSchema "pi-inventory/modules/stock/schema"
	warehouseSchema "pi-inventory/modules/warehouse/schema"

	commonModel "pi-inventory/common/models"

	"github.com/google/uuid"

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

	CreateDemoTaxList()
	CreateDemoSupplierList()
}

func CreateDemoTaxList() {
	for i := 0; i < 10; i++ {
		tax := models.TaxResponseBody{
			ID:          uuid.New(),
			TaxName:     "Demo Tax 1",
			AgencyID:    uuid.New(),
			Description: "Description",
			SalesRateHistory: []*models.RespRateHistory{
				{
					Rate:      10,
					StartDate: "",
				},
				{
					Rate:      20,
					StartDate: "",
				},
			},
			PurchaseRateHistory: []*models.RespRateHistory{
				{
					Rate:      5,
					StartDate: "",
				},
				{
					Rate:      6,
					StartDate: "",
				},
			},
			Status:    "active",
			CreatedBy: 1,
			AccountID: 1,
		}
		DemoTaxList = append(DemoTaxList, &tax)
	}
}

func CreateDemoSupplierList() {
	for i := 0; i < 10; i++ {
		supplier := models.SupplierResponseBody{
			ID:          uuid.New(),
			Title:       "Demo Title",
			FirstName:   "Demo First Name",
			LastName:    "Demo Last Name",
			DisplayName: "Supplier 1",
			Status:      "active",
		}
		DemoSupplierList = append(DemoSupplierList, &supplier)
	}
}
