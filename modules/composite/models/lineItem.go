package models

import (
	"pi-inventory/common/models"
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type LineItemResponseBody struct {
	Title         string            `json:"title"`
	StockID       uint64            `json:"stock_id"`
	AccountID     uint64            `json:"account_id"`
	Quantity      uint64            `json:"quantity"`
	UnitRate      decimal.Decimal   `json:"unit_rate"`
	SellingPrice  models.FloatMoney `json:"selling_price"`
	PurchaseRate  decimal.Decimal   `json:"purchase_rate"`
	PurchasePrice models.FloatMoney `json:"purchase_price"`
}

type AddLineItemRequestBody struct {
	Title         string          `json:"title"`
	StockID       uint64          `json:"stock_id"`
	AccountID     uint64          `json:"account_id"`
	Quantity      uint64          `json:"quantity"`
	UnitRate      decimal.Decimal `json:"unit_rate"`
	SellingPrice  models.Money    `json:"selling_price"`
	PurchaseRate  decimal.Decimal `json:"purchase_rate"`
	PurchasePrice models.Money    `json:"purchase_price"`
}

func validateLineItem(reqBody AddLineItemRequestBody) error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title,
			validation.Required,
			validation.Length(3, 100).ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RangeValidationErr,
				TranslationKey: "rangeValidationError",
				TranslationParams: map[string]interface{}{
					"field": "title",
					"min":   3,
					"max":   100,
				},
			}),
		),
		validation.Field(&reqBody.StockID,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Stock id",
				},
			}),
		),
		validation.Field(&reqBody.Quantity,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Quantity",
				},
			}),
		),
		validation.Field(&reqBody.UnitRate,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Unit rate",
				},
			}),
		),
		validation.Field(&reqBody.SellingPrice,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Selling price",
				},
			}),
			validation.By(models.ValidateMoney),
			validation.By(validateSellingPrice(reqBody.Quantity, reqBody.UnitRate, reqBody.SellingPrice)),
		),
		validation.Field(&reqBody.PurchaseRate,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Purchase rate",
				},
			}),
		),
		validation.Field(&reqBody.PurchasePrice,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Purchase price",
				},
			}),
			validation.By(models.ValidateMoney),
			validation.By(validatePurchasePrice(reqBody.Quantity, reqBody.PurchaseRate, reqBody.PurchasePrice)),
		),
	)
}

func validateSellingPrice(quantity uint64, unitRate decimal.Decimal, sellingPrice models.Money) validation.RuleFunc {
	return func(value interface{}) error {
		calculatedSellingPrice := unitRate.Mul(decimal.NewFromInt(int64(quantity)))
		if !calculatedSellingPrice.Equal(sellingPrice.Amount) {
			return &errors.ApplicationError{
				ErrorType:      errors.InvalidType,
				TranslationKey: "invalidLineItemPrice",
				TranslationParams: map[string]interface{}{
					"field": "Selling price",
					"rate":  "unit rate",
				},
			}
		}

		return nil
	}
}

func validatePurchasePrice(quantity uint64, purchaseRate decimal.Decimal, purchasePrice models.Money) validation.RuleFunc {
	return func(value interface{}) error {
		calculatedPurchasePrice := purchaseRate.Mul(decimal.NewFromInt(int64(quantity)))
		if !calculatedPurchasePrice.Equal(purchasePrice.Amount) {
			return &errors.ApplicationError{
				ErrorType:      errors.InvalidType,
				TranslationKey: "invalidLineItemPrice",
				TranslationParams: map[string]interface{}{
					"field": "Purchase price",
					"rate":  "purchase rate",
				},
			}
		}

		return nil
	}
}

type UpdateLineItemRequestBody struct {
}
