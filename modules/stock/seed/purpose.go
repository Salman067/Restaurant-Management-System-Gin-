package seed

import (
	"fmt"
	"pi-inventory/modules/stock/schema"
	"time"

	"gorm.io/gorm"

	purposeConst "pi-inventory/modules/stock/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func PurposeSeed(db *gorm.DB) {
	purposes := make([]*schema.Purpose, 0)
	for index := 0; index < 10; index++ {
		purpose := schema.Purpose{
			Title:     fakeData.BeerName(),
			OwnerID:   1,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		}
		purposes = append(purposes, &purpose)
	}
	result := db.Table(purposeConst.PurposeTable).Create(&purposes)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
