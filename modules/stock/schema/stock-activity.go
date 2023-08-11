package schema

import (
	"github.com/shopspring/decimal"
	"time"
)

type StockActivity struct {
	ID                            uint64 `gorm:"primaryKey"`
	OwnerID                       uint64
	AccountID                     uint64
	Mode                          string
	OperationType                 string
	LocationID                    uint64
	StockID                       uint64
	PurposeID                     uint64
	QuantityOnHand                uint64
	NewQuantity                   uint64
	AdjustedQuantity              uint64
	AdjustedDate                  *time.Time
	PurchaseDate                  *time.Time
	PurchasePreviousValue         decimal.Decimal
	PurchasePreviousValueCurrency string `gorm:"default:BDT"`
	PurchaseNewValue              decimal.Decimal
	PurchaseNewValueCurrency      string `gorm:"default:BDT"`
	PurchaseAdjustedValue         decimal.Decimal
	PurchaseAdjustedValueCurrency string `gorm:"default:BDT"`
	SellingPreviousValue          decimal.Decimal
	SellingPreviousValueCurrency  string `gorm:"default:BDT"`
	SellingNewValue               decimal.Decimal
	SellingNewValueCurrency       string `gorm:"default:BDT"`
	SellingAdjustedValue          decimal.Decimal
	SellingAdjustedValueCurrency  string `gorm:"default:BDT"`
	CreatedAt                     time.Time
	CreatedBy                     uint64
	UpdatedAt                     time.Time
	UpdatedBy                     uint64
}
