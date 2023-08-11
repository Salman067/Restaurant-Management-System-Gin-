package models

import (
	"pi-inventory/common/models"
	"pi-inventory/errors"
	attachmentModel "pi-inventory/modules/attachment/models"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type UpdateStockRequestBody struct {
	Name           string       `json:"name"`
	SKU            string       `json:"sku"`
	Status         string       `json:"status"`
	Description    string       `json:"description"`
	SellingPrice   models.Money `json:"selling_price"`
	PurchasePrice  models.Money `json:"purchase_price"`
	AsOfDate       *time.Time   `json:"as_of_date"`
	ExpiryDate     *time.Time   `json:"expiry_date"`
	PurchaseDate   *time.Time   `json:"purchase_date"`
	TrackInventory bool         `json:"track_inventory"`
	StockQty       uint64       `json:"stock_qty"`
	ReorderQty     uint64       `json:"reorder_qty"`
	AttachmentKey  string       `json:"attachment_key"`
	UnitID         uint64       `json:"unit_id"`
	CategoryID     uint64       `json:"category_id"`
	LocationID     uint64       `json:"location_id"`
	SupplierID     string       `json:"supplier_id"`
	SaleTaxID      string       `json:"sale_tax_id"`
	PurchaseTaxID  string       `json:"purchase_tax_id"`
	AccountID      uint64       `json:"account_id"`
}

func (reqBody UpdateStockRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Name, validation.Length(2, 100).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   2,
				"max":   100,
				"field": "name",
			},
		})),
		validation.Field(&reqBody.SKU, validation.Length(0, 100).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   0,
				"max":   100,
				"field": "sku",
			},
		})),
		/*validation.Field(&reqBody.StockQty, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "Stock quantity",
			},
		})),*/
		validation.Field(&reqBody.ReorderQty, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if _, ok := value.(uint64); !ok {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Reorder quantity",
						},
					}
				}
				return nil
			}),
		)),
		// validation.Field(&reqBody.SupplierID, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
		// 	ErrorType:      errors.RequiredField,
		// 	TranslationKey: "fieldRequired",
		// 	TranslationParams: map[string]interface{}{
		// 		"field": "Supplier id",
		// 	},
		// })),
		validation.Field(&reqBody.AsOfDate, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "As of date",
			},
		})),
		validation.Field(&reqBody.PurchaseDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.PurchaseDate != nil && reqBody.ExpiryDate != nil && reqBody.PurchaseDate.After(*reqBody.ExpiryDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanError",
						TranslationParams: map[string]interface{}{
							"field": "Expiry date",
							"value": "purchase date",
						},
					}
				}
				return nil
			}),
			validation.By(func(value interface{}) error {
				if reqBody.AsOfDate != nil && reqBody.PurchaseDate != nil && reqBody.PurchaseDate.After(*reqBody.AsOfDate) && !reqBody.PurchaseDate.Equal(*reqBody.AsOfDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "Purchase date",
							"field2": "as of date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.ExpiryDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.ExpiryDate != nil && reqBody.PurchaseDate != nil && reqBody.ExpiryDate.Before(*reqBody.PurchaseDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "greaterThanError",
						TranslationParams: map[string]interface{}{
							"field": "Expiry date",
							"value": "purchase date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.AsOfDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.AsOfDate != nil && reqBody.PurchaseDate != nil && reqBody.AsOfDate.Before(*reqBody.PurchaseDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "greaterThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "As of date",
							"field2": "purchase date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.SellingPrice,
			validation.Skip.When(reqBody.SellingPrice == models.Money{}),
			validation.By(models.ValidateMoney),
		),
		validation.Field(&reqBody.PurchasePrice,
			validation.Skip.When(reqBody.PurchasePrice == models.Money{}),
			validation.By(models.ValidateMoney),
		),
	)
}

type SingleStockResponse struct {
	ID            uint64                                      `json:"id"`
	OwnerID       uint64                                      `json:"owner_id"`
	AccountID     uint64                                      `json:"account_id"`
	Name          string                                      `json:"name"`
	Type          string                                      `json:"type"`
	SKU           string                                      `json:"sku"`
	Status        string                                      `json:"status"`
	Description   string                                      `json:"description"`
	SellingPrice  models.FloatMoney                           `json:"selling_price"`
	PurchasePrice models.FloatMoney                           `json:"purchase_price"`
	StockQty      uint64                                      `json:"stock_qty"`
	ReorderQty    uint64                                      `json:"reorder_qty"`
	RecordType    string                                      `json:"record_type"`
	AsOfDate      *time.Time                                  `json:"as_of_date"`
	PurchaseDate  *time.Time                                  `json:"purchase_date"`
	ExpiryDate    *time.Time                                  `json:"expiry_date"`
	Attachments   []*attachmentModel.AttachmentCustomResponse `json:"attachments"`
	SupplierID    string                                      `json:"supplier_id"`
	SaleTaxID     string                                      `json:"sale_tax_id"`
	PurchaseTaxID string                                      `json:"purchase_tax_id"`
	UnitTitle     string                                      `json:"unit_title"`
	CategoryTitle string                                      `json:"category_title"`
	LocationTitle string                                      `json:"location_title"`
}

type SingleStockResponseView struct {
	ID            uint64                                      `json:"id"`
	OwnerID       uint64                                      `json:"owner_id"`
	AccountID     uint64                                      `json:"account_id"`
	Name          string                                      `json:"name"`
	Type          string                                      `json:"type"`
	SKU           string                                      `json:"sku"`
	Status        string                                      `json:"status"`
	Description   string                                      `json:"description"`
	SellingPrice  models.FloatMoney                           `json:"selling_price"`
	PurchasePrice models.FloatMoney                           `json:"purchase_price"`
	StockQty      uint64                                      `json:"stock_qty"`
	ReorderQty    uint64                                      `json:"reorder_qty"`
	RecordType    string                                      `json:"record_type"`
	AsOfDate      *time.Time                                  `json:"as_of_date"`
	PurchaseDate  *time.Time                                  `json:"purchase_date"`
	ExpiryDate    *time.Time                                  `json:"expiry_date"`
	Attachments   []*attachmentModel.AttachmentCustomResponse `json:"attachments"`
	SupplierID    string                                      `json:"supplier_id"`
	SaleTaxID     string                                      `json:"sale_tax_id"`
	PurchaseTaxID string                                      `json:"purchase_tax_id"`
	UnitID        uint64                                      `json:"unit_id"`
	CategoryID    uint64                                      `json:"category_id"`
	LocationID    uint64                                      `json:"location_id"`
}

type AddStockRequestBody struct {
	Name           string       `json:"name"`
	SKU            string       `json:"sku"`
	Status         string       `json:"status"`
	Description    string       `json:"description"`
	SellingPrice   models.Money `json:"selling_price"`
	PurchasePrice  models.Money `json:"purchase_price"`
	AsOfDate       *time.Time   `json:"as_of_date"`
	ExpiryDate     *time.Time   `json:"expiry_date"`
	PurchaseDate   *time.Time   `json:"purchase_date"`
	TrackInventory bool         `json:"track_inventory"`
	StockQty       uint64       `json:"stock_qty"`
	ReorderQty     uint64       `json:"reorder_qty"`
	AttachmentKey  string       `json:"attachment_key"`
	UnitID         uint64       `json:"unit_id"`
	CategoryID     uint64       `json:"category_id"`
	LocationID     uint64       `json:"location_id"`
	SupplierID     string       `json:"supplier_id"`
	SaleTaxID      string       `json:"sale_tax_id"`
	PurchaseTaxID  string       `json:"purchase_tax_id"`
	GroupItemID    uint64       `json:"group_item_id"`
	AccountID      uint64       `json:"account_id"`
}

func (reqBody AddStockRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Name, validation.Required, validation.Length(2, 100).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   2,
				"max":   100,
				"field": "name",
			},
		})),
		validation.Field(&reqBody.SKU, validation.Length(0, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   0,
				"max":   50,
				"field": "sku",
			},
		})),
		// validation.Field(&reqBody.SupplierID, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
		// 	ErrorType:      errors.RequiredField,
		// 	TranslationKey: "fieldRequired",
		// 	TranslationParams: map[string]interface{}{
		// 		"field": "Supplier id",
		// 	},
		// })),
		validation.Field(&reqBody.StockQty, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "Stock quantity",
			},
		})),
		validation.Field(&reqBody.ReorderQty, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if _, ok := value.(uint64); !ok {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Reorder quantity",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.AsOfDate, validation.Required.When(reqBody.TrackInventory).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "As of date",
			},
		})),
		validation.Field(&reqBody.PurchaseDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.PurchaseDate != nil && reqBody.ExpiryDate != nil && reqBody.PurchaseDate.After(*reqBody.ExpiryDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanError",
						TranslationParams: map[string]interface{}{
							"field": "Expiry date",
							"value": "purchase date",
						},
					}
				}
				return nil
			}),
			validation.By(func(value interface{}) error {
				if reqBody.AsOfDate != nil && reqBody.PurchaseDate != nil && reqBody.PurchaseDate.After(*reqBody.AsOfDate) && !reqBody.PurchaseDate.Equal(*reqBody.AsOfDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "Purchase date",
							"field2": "as of date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.ExpiryDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.ExpiryDate != nil && reqBody.PurchaseDate != nil && reqBody.ExpiryDate.Before(*reqBody.PurchaseDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "greaterThanError",
						TranslationParams: map[string]interface{}{
							"field": "Expiry date",
							"value": "purchase date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.AsOfDate, validation.When(reqBody.TrackInventory,
			validation.By(func(value interface{}) error {
				if reqBody.AsOfDate != nil && reqBody.PurchaseDate != nil && reqBody.AsOfDate.Before(*reqBody.PurchaseDate) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "greaterThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "As of date",
							"field2": "purchase date",
						},
					}
				}
				return nil
			}),
		)),
		validation.Field(&reqBody.SellingPrice,
			validation.Skip.When(reqBody.SellingPrice == models.Money{}),
			validation.By(models.ValidateMoney),
			validation.By(func(value interface{}) error {
				return models.ValidateMoneyAmount(reqBody.SellingPrice, "amount")
			}),
		),
		validation.Field(&reqBody.PurchasePrice,
			validation.Skip.When(reqBody.PurchasePrice == models.Money{}),
			validation.By(models.ValidateMoney),
			validation.By(func(value interface{}) error {
				return models.ValidateMoneyAmount(reqBody.PurchasePrice, "amount")
			}),
		),
	)
}

type StockQueryParams struct {
	Type     string
	Status   string
	ListType string
	KeyWord  string
}

type StockSummaryResponse struct {
	CountOfStockList   int64                  `json:"count_of_stock_list"`
	CountOfItemList    int64                  `json:"count_of_item_list"`
	CountOfExpiryDate  int64                  `json:"count_of_expiry_date"`
	OutOfStock         int64                  `json:"out_of_stock"`
	LowOnStock         int64                  `json:"low_on_stock"`
	OutOfStockPageInfo *models.PageResponse   `json:"out_of_stock_page_info"`
	OutOfStocks        []*SingleStockResponse `json:"out_of_stocks"`
	LowOnStockPageInfo *models.PageResponse   `json:"low_on_stock_page_info"`
	LowOnStocks        []*SingleStockResponse `json:"low_on_stocks"`
}

type CustomQueryStockResponse struct {
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
	AttachmentKey         string
	UnitID                uint64
	CategoryID            uint64
	LocationID            uint64
	UnitTitle             string
	CategoryTitle         string
	LocationTitle         string
	SupplierID            string
	SaleTaxID             string
	PurchaseTaxID         string
	OwnerID               uint64
	AccountID             uint64
	CreatedAt             *time.Time
	CreatedBy             uint64
	UpdatedAt             *time.Time
	UpdatedBy             uint64
}

type RequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
}
