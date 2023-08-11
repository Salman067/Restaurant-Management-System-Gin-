package schema

import (
	"time"

	"github.com/shopspring/decimal"
)

type LineItem struct {
	ID                    uint64 `gorm:"primaryKey"`
	Title                 string
	OwnerID               uint64
	AccountID             uint64
	StockID               uint64
	Quantity              uint64
	UnitRate              decimal.Decimal
	SellingPriceCurrency  string `gorm:"default:BDT"`
	TotalSellingPrice     decimal.Decimal
	PurchaseRate          decimal.Decimal
	PurchasePriceCurrency string `gorm:"default:BDT"`
	TotalPurchasePrice    decimal.Decimal
	LineItemKey           string
	IsDeleted             bool
	CreatedAt             time.Time
	CreatedBy             uint64
	UpdatedAt             time.Time
	UpdatedBy             uint64
}
