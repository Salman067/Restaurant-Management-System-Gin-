package schema

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type GroupItem struct {
	ID                    uint64 `gorm:"primaryKey"`
	OwnerID               uint64
	AccountID             uint64
	Name                  string
	Tag                   string
	GroupItemUnit         string
	AttachmentKey         string
	Variant_values        datatypes.JSON
	CreatedAt             time.Time
	CreatedBy             uint64
	UpdatedAt             time.Time
	UpdatedBy             uint64
	SellingPrice          decimal.Decimal
	SellingPriceCurrency  string `gorm:"default:BDT"`
	PurchasePrice         decimal.Decimal
	PurchasePriceCurrency string `gorm:"default:BDT"`
}
