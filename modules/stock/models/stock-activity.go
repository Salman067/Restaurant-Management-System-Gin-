package models

import (
	"pi-inventory/common/models"
	"pi-inventory/errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type AddStockActivityRequestBody struct {
	Mode                  string       `json:"mode"`
	OperationType         string       `json:"operation_type"`
	StockID               uint64       `json:"stock_id"`
	AccountID             uint64       `json:"account_id"`
	PurposeID             uint64       `json:"purpose_id"`
	AdjustedDate          *time.Time   `json:"adjusted_date"`
	PurchaseDate          *time.Time   `json:"purchase_date"`
	LocationID            uint64       `json:"location_id"`
	QuantityOnHand        uint64       `json:"quantity_on_hand"`
	NewQuantity           uint64       `json:"new_quantity"`
	AdjustedQuantity      uint64       `json:"adjusted_quantity"`
	PreviousSellingPrice  models.Money `json:"previous_selling_price"`
	PreviousPurchasePrice models.Money `json:"previous_purchase_price"`
	NewSellingPrice       models.Money `json:"new_selling_price"`
	NewPurchasePrice      models.Money `json:"new_purchase_price"`
	AdjustedSellingPrice  models.Money `json:"adjusted_selling_price"`
	AdjustedPurchasePrice models.Money `json:"adjusted_purchase_price"`
}

func (reqBody AddStockActivityRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Mode,
			validation.By(validateMode), validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "modeRequired",
				TranslationParams: map[string]interface{}{
					"field": "Mode",
				},
			})),
		validation.Field(&reqBody.StockID,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Stock id",
				},
			})),
		validation.Field(&reqBody.PurposeID,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Purpose id",
				},
			})),
		validation.Field(&reqBody.OperationType,
			validation.By(validateOperationType), validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Operation type",
				},
			})),
		validation.Field(&reqBody.QuantityOnHand,
			validation.When(reqBody.Mode == "quantity",
				validation.By(func(value interface{}) error {
					if _, ok := value.(uint64); !ok {
						return &errors.ApplicationError{
							ErrorType:      errors.RequiredField,
							TranslationKey: "fieldRequired",
							TranslationParams: map[string]interface{}{
								"field": "Quantity on hand",
							},
						}
					}
					return nil
				}),
			),
		),
		validation.Field(&reqBody.NewQuantity,
			validation.Required.When(reqBody.Mode == "quantity").ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "New quantity",
				},
			}),
			validation.By(func(value interface{}) error {
				if reqBody.Mode == "quantity" && reqBody.OperationType == "sub" && reqBody.NewQuantity > reqBody.QuantityOnHand {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "New quantity",
							"field2": "quantity on hand",
						},
					}
				}
				return nil
			}),
		),
		validation.Field(&reqBody.AdjustedQuantity,
			validation.When(reqBody.Mode == "quantity",
				validation.By(func(value interface{}) error {
					if _, ok := value.(uint64); !ok {
						return &errors.ApplicationError{
							ErrorType:      errors.RequiredField,
							TranslationKey: "fieldRequired",
							TranslationParams: map[string]interface{}{
								"field": "Adjusted quantity",
							},
						}
					}
					return nil
				}),
			),
			validation.By(func(value interface{}) error {
				result := reqBody.QuantityOnHand - reqBody.NewQuantity
				if reqBody.Mode == "quantity" && reqBody.OperationType == "sub" && reqBody.AdjustedQuantity != result {
					return &errors.ApplicationError{
						ErrorType:      errors.Unsuccessfull,
						TranslationKey: "SubtractionError",
						TranslationParams: map[string]interface{}{
							"field": "Adjusted quantity",
						},
					}
				}
				return nil
			}),
			validation.By(func(value interface{}) error {
				result := reqBody.QuantityOnHand + reqBody.NewQuantity
				if reqBody.Mode == "quantity" && reqBody.OperationType == "add" && reqBody.AdjustedQuantity != result {
					return &errors.ApplicationError{
						ErrorType:      errors.Unsuccessfull,
						TranslationKey: "AdditionError",
						TranslationParams: map[string]interface{}{
							"field": "Adjusted quantity",
						},
					}
				}
				return nil
			}),
		),
		validation.Field(&reqBody.PreviousSellingPrice, validation.When(reqBody.Mode == "value", validation.By(models.ValidateMoney))),
		validation.Field(&reqBody.NewSellingPrice, validation.When(reqBody.Mode == "value", validation.By(models.ValidateNewSellingPrice)),
			validation.By(func(value interface{}) error {
				if reqBody.OperationType == "sub" && reqBody.NewSellingPrice.Amount.GreaterThan(reqBody.PreviousSellingPrice.Amount) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "lessThanOrEqualError",
						TranslationParams: map[string]interface{}{
							"field1": "New selling price",
							"field2": "previous selling price",
						},
					}
				}
				return nil
			}),
		),
		validation.Field(&reqBody.AdjustedSellingPrice,
			validation.When(reqBody.Mode == "value", validation.By(models.ValidateMoney)),
			validation.By(func(value interface{}) error {
				result := reqBody.PreviousSellingPrice.Amount.Sub(reqBody.NewSellingPrice.Amount)
				if reqBody.Mode == "value" && reqBody.OperationType == "sub" && (reqBody.AdjustedSellingPrice.Amount.GreaterThan(result) || reqBody.AdjustedSellingPrice.Amount.LessThan(result)) {
					return &errors.ApplicationError{
						ErrorType:      errors.Unsuccessfull,
						TranslationKey: "SubtractionError",
						TranslationParams: map[string]interface{}{
							"field": "Adjusted selling price",
						},
					}
				}
				return nil
			}),
			validation.By(func(value interface{}) error {
				result := reqBody.PreviousSellingPrice.Amount.Add(reqBody.NewSellingPrice.Amount)
				if reqBody.Mode == "value" && reqBody.OperationType == "add" && (reqBody.AdjustedSellingPrice.Amount.GreaterThan(result) || reqBody.AdjustedSellingPrice.Amount.LessThan(result)) {
					return &errors.ApplicationError{
						ErrorType:      errors.Unsuccessfull,
						TranslationKey: "AdditionError",
						TranslationParams: map[string]interface{}{
							"field": "Adjusted selling price",
						},
					}
				}
				return nil
			}),
		),
	)
}

func validateMode(value interface{}) error {
	mode, ok := value.(string)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "mode",
			},
		}
	}

	if mode != "quantity" && mode != "value" {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "modeValueError",
			TranslationParams: map[string]interface{}{
				"field": "mode",
			},
		}
	}

	return nil
}
func validateOperationType(value interface{}) error {
	mode, ok := value.(string)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "operation type",
			},
		}
	}

	if mode != "add" && mode != "sub" {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "operationTypeValueError",
			TranslationParams: map[string]interface{}{
				"field": "operation type",
			},
		}
	}

	return nil
}

type SingleStockActivityResponse struct {
	ID                   uint64            `json:"id"`
	OwnerID              uint64            `json:"owner_id"`
	AccountID            uint64            `json:"account_id"`
	Mode                 string            `json:"mode"`
	OperationType        string            `json:"operation_type"`
	StockID              uint64            `json:"stock_id"`
	StockTitle           string            `json:"stock_title"`
	PurposeID            uint64            `json:"purpose_id"`
	PurposeTitle         string            `json:"purpose_title"`
	NewQuantity          uint64            `json:"new_quantity"`
	AdjustedQuantity     uint64            `json:"adjusted_quantity"`
	AdjustedDate         *time.Time        `json:"adjusted_date"`
	PurchaseDate         *time.Time        `json:"purchase_date"`
	NewSellingPrice      models.FloatMoney `json:"new_selling_price"`
	AdjustedSellingPrice models.FloatMoney `json:"adjusted_selling_price"`
}

type StockActivityQueryParams struct {
	StockID uint64
	Keyword string
}

type CustomStockActivityResponse struct {
	ID                            uint64
	OwnerID                       uint64
	AccountID                     uint64
	Mode                          string
	OperationType                 string
	LocationID                    uint64
	StockID                       uint64
	StockTitle                    string
	PurposeID                     uint64
	PurposeTitle                  string
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
type BulkAdjustment struct {
	StockActivities []*AddStockActivityRequestBody `json:"stock_activities"`
}

func (ba *BulkAdjustment) Validate() error {
	return validation.ValidateStruct(ba,
		validation.Field(&ba.StockActivities,
			validation.Each(validation.By(validateBulkAdjustmentWrapper))))

}

func validateBulkAdjustmentWrapper(value interface{}) error {
	var stockActivity, _ = value.(AddStockActivityRequestBody)
	return stockActivity.Validate()
}
