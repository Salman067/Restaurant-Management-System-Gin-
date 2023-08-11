package seed

import (
	"fmt"
	stockActivityConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/schema"
	"time"

	fakeData "github.com/brianvoe/gofakeit/v6"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func StockActivitySeed(db *gorm.DB) {
	stockActivities := make([]*schema.StockActivity, 0)
	for index := 0; index < 10; index++ {
		stockActivity := schema.StockActivity{
			Mode:                          fakeData.CarModel(),
			OperationType:                 fakeData.Adverb(),
			StockID:                       2,
			PurposeID:                     1,
			LocationID:                    2,
			QuantityOnHand:                5,
			NewQuantity:                   8,
			AdjustedQuantity:              3,
			PurchasePreviousValue:         decimal.New(11, 1),
			PurchasePreviousValueCurrency: fakeData.CurrencyShort(),
			SellingPreviousValue:          decimal.New(11, 1),
			SellingPreviousValueCurrency:  fakeData.CurrencyShort(),
			PurchaseNewValue:              decimal.New(11, 1),
			PurchaseNewValueCurrency:      fakeData.CurrencyShort(),
			SellingNewValue:               decimal.New(11, 1),
			SellingNewValueCurrency:       fakeData.CurrencyShort(),
			PurchaseAdjustedValue:         decimal.New(11, 1),
			PurchaseAdjustedValueCurrency: fakeData.CurrencyShort(),
			SellingAdjustedValue:          decimal.New(11, 1),
			SellingAdjustedValueCurrency:  fakeData.CurrencyShort(),
			CreatedAt:                     time.Now(),
			CreatedBy:                     1,
		}
		stockActivities = append(stockActivities, &stockActivity)
	}
	result := db.Table(stockActivityConst.StockActivityTable).Create(&stockActivities)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
