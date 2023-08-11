package seed

import (
	"fmt"
	"gorm.io/gorm"
	"pi-inventory/modules/warehouse/schema"
	"time"

	warehouseConst "pi-inventory/modules/warehouse/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func WarehouseSeed(db *gorm.DB) {

	categories := make([]*schema.Warehouse, 0)

	for index := 0; index < 10; index++ {
		Warehouse := schema.Warehouse{
			Title:     fakeData.Company(),
			Address:   fakeData.Address().City,
			OwnerID:   1,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		}
		categories = append(categories, &Warehouse)
	}
	result := db.Table(warehouseConst.WarehouseTable).Create(&categories)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
