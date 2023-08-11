package models

import (
	"pi-inventory/common/models"
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
	attachmentModel "pi-inventory/modules/attachment/models"
)

type UpdateCompositeRequestBody struct {
	Title         string                    `json:"title"`
	Description   string                    `json:"description"`
	AccountID     uint64                    `json:"account_id"`
	SellingPrice  models.Money              `json:"selling_price"`
	PurchasePrice models.Money              `json:"purchase_price"`
	LineItemKey   string                    `json:"line_item_key"`
	LineItems     *[]AddLineItemRequestBody `json:"line_items"`
}

func (reqBody UpdateCompositeRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Length(3, 100).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"field": "title",
				"min":   "3",
				"max":   "100",
			},
		})),
	)
}

type AddCompositeRequestBody struct {
	Title         string                    `json:"title"`
	Description   string                    `json:"description"`
	Tag           string                    `json:"tag"`
	AccountID     uint64                    `json:"account_id"`
	SellingPrice  models.Money              `json:"selling_price"`
	PurchasePrice models.Money              `json:"purchase_price"`
	AttachmentKey string                    `json:"attachment_key"`
	LineItems     *[]AddLineItemRequestBody `json:"line_items"`
}

func (reqBody AddCompositeRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "title",
				},
			}),
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
		validation.Field(&reqBody.Tag,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "tag",
				},
			}),
			validation.Length(3, 100).ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RangeValidationErr,
				TranslationKey: "rangeValidationError",
				TranslationParams: map[string]interface{}{
					"field": "tag",
					"min":   3,
					"max":   100,
				},
			}),
		),
		validation.Field(&reqBody.LineItems,
			validation.By(func(value interface{}) error {
				if isLineItemListEmpty(reqBody.LineItems, nil) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Line items",
						},
					}
				}
				return nil
			}),
			validation.By(validateAllLineItems),
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
			validation.By(func(value interface{}) error {
				return validateCompositePrices(reqBody.LineItems, "SellingPrice", reqBody.SellingPrice.Amount)
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
			validation.By(func(value interface{}) error {
				return validateCompositePrices(reqBody.LineItems, "PurchasePrice", reqBody.PurchasePrice.Amount)
			}),
		),
	)
}
func isLineItemListEmpty(lineItemWhileCreate *[]AddLineItemRequestBody, lineItemWhileUpdating *[]AddLineItemRequestBody) bool {
	if lineItemWhileCreate != nil {
		lineItem := *lineItemWhileCreate
		if len(lineItem) == 0 {
			return true
		}
	}
	if lineItemWhileUpdating != nil {
		lineItem := *lineItemWhileUpdating
		if len(lineItem) == 0 {
			return true
		}
	}
	return false
}
func validateAllLineItems(items interface{}) error {
	lineItems, ok := items.(*[]AddLineItemRequestBody)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidLineItemsValue",
			TranslationParams: map[string]interface{}{
				"field": "Line items",
			},
		}
	}

	for _, item := range *lineItems {
		err := validateLineItem(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateCompositePrices(items interface{}, targetPrice string, PriceToCheck decimal.Decimal) error {
	lineItems, ok := items.(*[]AddLineItemRequestBody)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "line items",
			},
		}
	}
	var totalSellingPrice decimal.Decimal
	var totalPurchasePrice decimal.Decimal

	for _, item := range *lineItems {
		totalSellingPrice = totalSellingPrice.Add(item.SellingPrice.Amount)
		totalPurchasePrice = totalPurchasePrice.Add(item.PurchasePrice.Amount)
	}
	if targetPrice == "SellingPrice" && totalSellingPrice.Cmp(PriceToCheck) != 0 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidPriceValue",
			TranslationParams: map[string]interface{}{
				"field": "Selling price",
			},
		}
	} else if targetPrice == "PurchasePrice" && totalPurchasePrice.Cmp(PriceToCheck) != 0 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidPriceValue",
			TranslationParams: map[string]interface{}{
				"field": "Purchase price",
			},
		}
	}
	return nil
}

type SingleCompositeResponseBody struct {
	ID            uint64                                      `json:"id"`
	OwnerID       uint64                                      `json:"owner_id"`
	AccountID     uint64                                      `json:"account_id"`
	Title         string                                      `json:"title"`
	Description   string                                      `json:"description"`
	Tag           string                                      `json:"tag"`
	SellingPrice  models.FloatMoney                           `json:"selling_price"`
	PurchasePrice models.FloatMoney                           `json:"purchase_price"`
	AttachmentKey string                                      `json:"attachment_key"`
	Attachments   []*attachmentModel.AttachmentCustomResponse `json:"attachments"`
	LineItems     []*LineItemResponseBody                     `json:"line_items"`
}

type CompositeQueryParams struct {
	KeyWord string
}

type RequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
}
