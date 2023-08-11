package schema

import (
	"time"

	"github.com/shopspring/decimal"
)

type Composite struct {
	ID                    uint64 `gorm:"primaryKey"`
	Title                 string
	OwnerID               uint64
	AccountID             uint64
	Tag                   string
	Description           string
	SellingPrice          decimal.Decimal
	SellingPriceCurrency  string `gorm:"default:BDT"`
	PurchasePrice         decimal.Decimal
	PurchasePriceCurrency string `gorm:"default:BDT"`
	AttachmentKey         string
	LineItemKey           string
	CreatedAt             time.Time
	CreatedBy             uint64
	UpdatedAt             time.Time
	UpdatedBy             uint64
}
