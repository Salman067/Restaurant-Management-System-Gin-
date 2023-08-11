package seed

import (
	"fmt"
	"gorm.io/gorm"
	"pi-inventory/modules/stock/schema"
	"time"

	unitConst "pi-inventory/modules/stock/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func UnitSeed(db *gorm.DB) {
	units := make([]*schema.Unit, 0)
	for index := 0; index < 10; index++ {
		unit := schema.Unit{
			Title:     fakeData.BeerName(),
			OwnerID:   1,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		}
		units = append(units, &unit)
	}
	result := db.Table(unitConst.UnitTable).Create(&units)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
