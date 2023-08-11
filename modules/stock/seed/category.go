package seed

import (
	"fmt"
	"gorm.io/gorm"
	"pi-inventory/modules/stock/schema"
	"time"

	stockConst "pi-inventory/modules/stock/consts"

	fakeData "github.com/brianvoe/gofakeit/v6"
)

func CategorySeed(db *gorm.DB) {

	categories := make([]*schema.Category, 0)

	for index := 0; index < 10; index++ {
		category := schema.Category{
			Title:     fakeData.BeerName(),
			OwnerID:   1,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		}
		categories = append(categories, &category)
	}
	result := db.Table(stockConst.CategoryTable).Create(&categories)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
