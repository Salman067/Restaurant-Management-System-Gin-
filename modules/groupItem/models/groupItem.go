package models

import (
	"pi-inventory/common/models"
	"pi-inventory/errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	attachmentModel "pi-inventory/modules/attachment/models"
)

type RequestGroupItemBody struct {
	Name              string                  `json:"name"`
	Tag               string                  `json:"tag"`
	GroupItemUnit     string                  `json:"group_item_unit"`
	SellingPriceCost  models.Money            `json:"selling_price_cost"`
	PurchasePriceCost models.Money            `json:"purchase_price_cost"`
	AttachmentKey     string                  `json:"attachment_key"`
	Variants          []*RequestVariant       `json:"variants"`
	GroupLineItems    []*RequestGroupLineItem `json:"group_line_items"`
	AccountID         uint64                  `json:"account_id"`
}

func (reqBody RequestGroupItemBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Name,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "name",
				},
			}),
			validation.Length(3, 100).ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RangeValidationErr,
				TranslationKey: "rangeValidationError",
				TranslationParams: map[string]interface{}{
					"field": "name",
					"min":   3,
					"max":   100,
				},
			}),
		),
		validation.Field(&reqBody.Variants,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "variants",
				},
			}),
		),
		validation.Field(&reqBody.GroupLineItems,
			validation.By(func(value interface{}) error {
				if isGroupLineItemListEmpty(reqBody.GroupLineItems, nil) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "group line items",
						},
					}
				}
				return nil
			}),
			validation.By(validateAllGroupLineItems),
		),
		validation.Field(&reqBody.Variants,
			validation.By(func(value interface{}) error {
				if isRequestVariantsEmpty(reqBody.Variants, nil) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "variants",
						},
					}
				}
				return nil
			}),
			validation.By(validateAllVariants),
		),
		// validation.Field(&reqBody.SellingPriceCost,
		// 	validation.Required.ErrorObject(&errors.ApplicationError{
		// 		ErrorType:      errors.RequiredField,
		// 		TranslationKey: "fieldRequired",
		// 		TranslationParams: map[string]interface{}{
		// 			"field": "selling_price",
		// 		},
		// 	}),
		// 	validation.By(func(value interface{}) error {
		// 		if reqBody.SellingPriceCost.Amount.LessThan(decimal.Zero) {
		// 			return &errors.ApplicationError{
		// 				ErrorType:      errors.RequiredField,
		// 				TranslationKey: "fieldRequired",
		// 				TranslationParams: map[string]interface{}{
		// 					"field": "selling_price_amount",
		// 				},
		// 			}
		// 		}
		// 		return models.ValidateMoneyAmount(reqBody.SellingPriceCost, "amount")
		// 	}),
		// validation.By(func(value interface{}) error {
		// 	currency := reqBody.SellingPriceCost.Currency
		// 	if len(currency) == 0 {
		// 		return &errors.ApplicationError{
		// 			ErrorType:      errors.RequiredField,
		// 			TranslationKey: "fieldRequired",
		// 			TranslationParams: map[string]interface{}{
		// 				"field": "selling_price_currency",
		// 			},
		// 		}
		// 	}
		// 	return nil
		// }),
		// validation.By(func(value interface{}) error {
		// 	return validateGroupItemPrices(reqBody.GroupLineItems, "SellingPriceCost", reqBody.SellingPriceCost.Amount)
		// }),
		// ),
		// validation.Field(&reqBody.PurchasePriceCost,
		// 	validation.Required.ErrorObject(&errors.ApplicationError{
		// 		ErrorType:      errors.RequiredField,
		// 		TranslationKey: "fieldRequired",
		// 		TranslationParams: map[string]interface{}{
		// 			"field": "purchase_price_cost",
		// 		},
		// 	}),
		// 	validation.By(func(value interface{}) error {
		// 		if reqBody.PurchasePriceCost.Amount.LessThan(decimal.Zero) {
		// 			return &errors.ApplicationError{
		// 				ErrorType:      errors.RequiredField,
		// 				TranslationKey: "fieldRequired",
		// 				TranslationParams: map[string]interface{}{
		// 					"field": "purchase_price_amount",
		// 				},
		// 			}
		// 		}
		// 		return models.ValidateMoneyAmount(reqBody.PurchasePriceCost, "amount")
		// 	}),
		// 	validation.By(func(value interface{}) error {
		// 		currency := reqBody.PurchasePriceCost.Currency
		// 		if len(currency) == 0 {
		// 			return &errors.ApplicationError{
		// 				ErrorType:      errors.RequiredField,
		// 				TranslationKey: "fieldRequired",
		// 				TranslationParams: map[string]interface{}{
		// 					"field": "purchase_price_currency",
		// 				},
		// 			}
		// 		}
		// 		return nil
		// 	}),
		// 	validation.By(func(value interface{}) error {
		// 		return validateGroupItemPrices(reqBody.GroupLineItems, "PurchasePriceCost", reqBody.PurchasePriceCost.Amount)
		// 	}),
		// ),
	)
}
func isGroupLineItemListEmpty(lineItemWhileCreate []*RequestGroupLineItem, lineItemWhileUpdating []*RequestGroupLineItem) bool {
	if lineItemWhileCreate != nil {
		lineItem := lineItemWhileCreate
		if len(lineItem) == 0 {
			return true
		}
	}
	if lineItemWhileUpdating != nil {
		lineItem := lineItemWhileUpdating
		if len(lineItem) == 0 {
			return true
		}
	}
	return false
}
func validateAllGroupLineItems(items interface{}) error {
	lineItems, ok := items.([]*RequestGroupLineItem)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidGroupLineItemsValue",
			TranslationParams: map[string]interface{}{
				"field": "group line items",
			},
		}
	}

	for _, item := range lineItems {
		err := validateGroupLineItem(*item)
		if err != nil {
			return err
		}
	}

	return nil
}

func isRequestVariantsEmpty(variantsWhileCreate []*RequestVariant, variantsWhileUpdating []*RequestVariant) bool {
	if variantsWhileCreate != nil {
		variants := variantsWhileCreate
		if len(variants) == 0 {
			return true
		}
	}
	if variantsWhileUpdating != nil {
		variants := variantsWhileUpdating
		if len(variants) == 0 {
			return true
		}
	}
	return false
}

func validateAllVariants(items interface{}) error {
	variants, ok := items.([]*RequestVariant)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "variant",
			},
		}
	}

	for _, variant := range variants {
		err := validateRequestVariant(*variant)
		if err != nil {
			return err
		}
	}

	return nil
}

// func validateGroupItemPrices(items interface{}, targetPrice string, PriceToCheck decimal.Decimal) error {
// 	lineItems, ok := items.([]*RequestGroupLineItem)
// 	if !ok {
// 		return &errors.ApplicationError{
// 			ErrorType:      errors.RequiredField,
// 			TranslationKey: "invalidGroupLineItemsValue",
// 			TranslationParams: map[string]interface{}{
// 				"field": "group_line_items",
// 			},
// 		}
// 	}
// 	var totalSellingPrice decimal.Decimal
// 	var totalCostPrice decimal.Decimal

// 	for _, item := range lineItems {
// 		totalSellingPrice = totalSellingPrice.Add(item.SellingPrice.Amount)
// 		totalCostPrice = totalCostPrice.Add(item.CostPrice.Amount)
// 	}
// 	if targetPrice == "SellingPriceCost" && totalSellingPrice.Cmp(PriceToCheck) != 0 {
// 		return &errors.ApplicationError{
// 			ErrorType:      errors.RequiredField,
// 			TranslationKey: "invalidSellingPriceValue",
// 			TranslationParams: map[string]interface{}{
// 				"field": "groupItem_price",
// 			},
// 		}
// 	} else if targetPrice == "PurchasePriceCost" && totalCostPrice.Cmp(PriceToCheck) != 0 {
// 		return &errors.ApplicationError{
// 			ErrorType:      errors.RequiredField,
// 			TranslationKey: "invalidPurchasePriceCostValue",
// 			TranslationParams: map[string]interface{}{
// 				"field": "groupItem_price",
// 			},
// 		}
// 	}
// 	return nil
// }

type ResponseGroupItemBody struct {
	ID                uint64                                      `json:"id"`
	Name              string                                      `json:"name"`
	Tag               string                                      `json:"tag"`
	AccountID         uint64                                      `json:"account_id"`
	GroupItemUnit     string                                      `json:"group_item_unit"`
	SellingPriceCost  models.FloatMoney                           `json:"selling_price_cost"`
	PurchasePriceCost models.FloatMoney                           `json:"purchase_price_cost"`
	Variants          datatypes.JSON                              `json:"variants"`
	Attachments       []*attachmentModel.AttachmentCustomResponse `json:"attachments"`
	GroupLineItems    []*ResponseGroupLineItem                    `json:"group_line_items"`
}
type RequestVariant struct {
	ID     uint64 `json:"variant_id"`
	Title  string `json:"variant_title"`
	Values string `json:"variant_values"`
}

func validateRequestVariant(reqBody RequestVariant) error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.ID,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "ID",
				},
			}),
		),
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
					"min":   2,
					"max":   50,
				},
			}),
		),
		validation.Field(&reqBody.Values,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "values",
				},
			}),
		),
	)
}

type SingleGroupItemResponse struct {
	ID                    uint64            `json:"id"`
	Name                  string            `json:"name"`
	AccountID             uint64            `json:"account_id"`
	Tag                   string            `json:"tag"`
	CountOfGroupLineItems int               `json:"count_of_group_line_items"`
	SellingPriceCost      models.FloatMoney `json:"selling_price_cost"`
	PurchasePriceCost     models.FloatMoney `json:"purchase_price_cost"`
}
type UpdateGroupItemRequestBody struct {
	Name              string                            `json:"name"`
	Tag               string                            `json:"tag"`
	GroupItemUnit     string                            `json:"group_item_unit"`
	AccountID         uint64                            `json:"account_id"`
	SellingPriceCost  models.Money                      `json:"selling_price_cost"`
	PurchasePriceCost models.Money                      `json:"purchase_price_cost"`
	AttachmentKey     string                            `json:"attachment_key"`
	Variants          []*RequestVariant                 `json:"variants"`
	GroupLineItems    []*UpdateGroupLineItemRequestBody `json:"group_line_items"`
}
type UpdateGroupLineItemRequestBody struct {
	ID            uint64       `json:"id"`
	Title         string       `json:"title"`
	AccountID     uint64       `json:"account_id"`
	SKU           string       `json:"sku"`
	SellingPrice  models.Money `json:"selling_price"`
	PurchasePrice models.Money `json:"purchase_price"`
	IsStocked     bool         `json:"is_stocked"`
	AsOfDate      *time.Time   `json:"as_of_date"`
	ExpiryDate    *time.Time   `json:"expiry_date"`
	PurchaseDate  *time.Time   `json:"purchase_date"`
	StockQty      uint64       `json:"stock_qty"`
	ReorderQty    uint64       `json:"reorder_qty"`
	UnitID        uint64       `json:"unit_id"`
	CategoryID    uint64       `json:"category_id"`
	LocationID    uint64       `json:"location_id"`
	SupplierID    string       `json:"supplier_id"`
}

type GroupItemQueryParams struct {
	KeyWord string
}

type RequestGroupLineItem struct {
	Title        string       `json:"title"`
	SKU          string       `json:"sku"`
	StockQty     uint64       `json:"stock_qty"`
	ReorderQty   uint64       `json:"reorder_qty"`
	SupplierID   string       `json:"supplier_id"`
	AccountID    uint64       `json:"account_id"`
	LocationID   uint64       `json:"location_id"`
	CostPrice    models.Money `json:"cost_price"`
	SellingPrice models.Money `json:"selling_price"`
	AsOfDate     *time.Time   `json:"as_of_date"`
	PurchaseDate *time.Time   `json:"purchase_date"`
	ExpiryDate   *time.Time   `json:"expiry_date"`
	IsStocked    bool         `json:"is_stocked"`
}

func validateGroupLineItem(reqBody RequestGroupLineItem) error {
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
		validation.Field(&reqBody.SKU,
			validation.Length(0, 50).ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RangeValidationErr,
				TranslationKey: "rangeValidationError",
				TranslationParams: map[string]interface{}{
					"field": "sku",
					"min":   0,
					"max":   50,
				},
			}),
		),
		//validation.Field(&reqBody.SupplierID, validation.Required.When(reqBody.IsStocked).ErrorObject(&errors.ApplicationError{
		//	ErrorType:      errors.RequiredField,
		//	TranslationKey: "fieldRequired",
		//	TranslationParams: map[string]interface{}{
		//		"field": "Supplier id",
		//	},
		//})),
		validation.Field(&reqBody.StockQty, validation.Required.When(reqBody.IsStocked).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "Stock quantity",
			},
		})),
		validation.Field(&reqBody.ReorderQty, validation.Required.When(reqBody.IsStocked).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "fieldRequired",
			TranslationParams: map[string]interface{}{
				"field": "Reorder quantity",
			},
		})),
		validation.Field(&reqBody.SellingPrice,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Selling price",
				},
			}),
			validation.By(func(value interface{}) error {
				if reqBody.SellingPrice.Amount.LessThan(decimal.Zero) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Selling price amount",
						},
					}
				}
				return models.ValidateMoneyAmount(reqBody.SellingPrice, "amount")
			}),
			validation.By(func(value interface{}) error {
				currency := reqBody.SellingPrice.Currency
				if len(currency) == 0 {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Selling price currency",
						},
					}
				}
				return nil
			}),
		),
		validation.Field(&reqBody.CostPrice,
			validation.Required.ErrorObject(&errors.ApplicationError{
				ErrorType:      errors.RequiredField,
				TranslationKey: "fieldRequired",
				TranslationParams: map[string]interface{}{
					"field": "Cost price",
				},
			}),
			validation.By(func(value interface{}) error {
				if reqBody.CostPrice.Amount.LessThan(decimal.Zero) {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Cost price amount",
						},
					}
				}
				return models.ValidateMoneyAmount(reqBody.CostPrice, "amount")
			}),
			validation.By(func(value interface{}) error {
				currency := reqBody.CostPrice.Currency
				if len(currency) == 0 {
					return &errors.ApplicationError{
						ErrorType:      errors.RequiredField,
						TranslationKey: "fieldRequired",
						TranslationParams: map[string]interface{}{
							"field": "Cost price currency",
						},
					}
				}
				return nil
			}),
		),
	)
}

type ResponseGroupLineItem struct {
	Title        string            `json:"title"`
	SKU          string            `json:"sku"`
	StockQty     uint64            `json:"stock_qty"`
	ReorderQty   uint64            `json:"reorder_qty"`
	StockID      uint64            `json:"stock_id"`
	AccountID    uint64            `json:"account_id"`
	SupplierID   string            `json:"supplier_id"`
	LocationID   uint64            `json:"location_id"`
	CostPrice    models.FloatMoney `json:"cost_price"`
	SellingPrice models.FloatMoney `json:"selling_price"`
	AsOfDate     *time.Time        `json:"as_of_date"`
	PurchaseDate *time.Time        `json:"purchase_date"`
	ExpiryDate   *time.Time        `json:"expiry_date"`
	IsStocked    bool              `json:"is_stocked"`
}

type RequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
}
