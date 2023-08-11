package seed

import (
	"fmt"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/schema"
	"time"

	fakeData "github.com/brianvoe/gofakeit/v6"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func StockSeed(db *gorm.DB) {
	stocks := make([]*schema.Stock, 0)
	for index := 0; index < 10; index++ {
		stock := schema.Stock{
			Name:                  fakeData.Name(),
			SKU:                   fakeData.RandomString(fakeData.NiceColors()),
			Type:                  "non-inventory",
			Status:                "active",
			Description:           fakeData.EmojiDescription(),
			SellingPrice:          decimal.New(10, 1),
			SellingPriceCurrency:  fakeData.CurrencyShort(),
			PurchasePrice:         decimal.New(11, 1),
			PurchasePriceCurrency: fakeData.CurrencyShort(),
			RecordType:            "stock",
			StockQty:              0,
			ReorderQty:            0,
			OwnerID:               1,
			UnitID:                2,
			CategoryID:            3,
			LocationID:            4,
			SupplierID:            "abc-def-ghi-jkl-mno",
			AttachmentKey:         fakeData.UUID(),
			CreatedAt:             time.Now(),
			CreatedBy:             1,
		}
		stocks = append(stocks, &stock)
	}
	result := db.Table(stockConst.StockTable).Create(&stocks)
	if result.Error != nil {
		panic(fmt.Errorf("failed generating statement: %w", result.Error))
	}
}
