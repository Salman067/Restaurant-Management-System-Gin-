package seed

import (
	"fmt"
	"pi-inventory/modules/groupItem/schema"
	"time"

	"gorm.io/gorm"

	groupItemConst "pi-inventory/modules/groupItem/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func VariantSeed(db *gorm.DB) {

	variants := make([]*schema.Variant, 0)

	for index := 0; index < 10; index++ {
		variant := schema.Variant{
			Title:     fakeData.BeerName(),
			OwnerID:   1,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		}
		variants = append(variants, &variant)
	}
	result := db.Table(groupItemConst.VariantTable).Create(&variants)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
