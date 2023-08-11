package schema

import (
	"github.com/google/uuid"
	"time"

	"github.com/shopspring/decimal"
)

type Stock struct {
	ID                    uint64 `gorm:"primaryKey"`
	Name                  string
	SKU                   string
	Type                  string
	Status                string `gorm:"default:active"`
	Description           string
	SellingPrice          decimal.Decimal
	SellingPriceCurrency  string `gorm:"default:BDT"`
	PurchasePrice         decimal.Decimal
	PurchasePriceCurrency string `gorm:"default:BDT"`
	TrackInventory        bool
	RecordType            string
	StockQty              uint64
	ReorderQty            uint64
	AsOfDate              *time.Time `gorm:"default:NULL"`
	PurchaseDate          *time.Time `gorm:"default:NULL"`
	ExpiryDate            *time.Time `gorm:"default:NULL"`
	UnitID                uint64
	AttachmentKey         string
	CategoryID            uint64
	LocationID            uint64
	SupplierID            string
	SaleTaxID             string
	PurchaseTaxID         string
	OwnerID               uint64
	AccountID             uint64
	GroupItemID           uint64
	StockUUID             uuid.UUID
	IsDeleted             bool
	CreatedAt             time.Time
	CreatedBy             uint64
	UpdatedAt             time.Time
	UpdatedBy             uint64
}
