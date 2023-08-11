package models

import (
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateWarehouseRequestBody struct {
	Title     string `json:"title"`
	Address   string `json:"address"`
	AccountID uint64 `json:"account_id"`
	IsUsed    bool   `json:"is_used"`
}

func (reqBody UpdateWarehouseRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Length(3, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   3,
				"max":   50,
				"field": "title",
			},
		})),
	)
}

type AddWarehouseRequestBody struct {
	Title     string `json:"title"`
	Address   string `json:"address"`
	AccountID uint64 `json:"account_id"`
	IsUsed    bool   `json:"is_used"`
}

func (reqBody AddWarehouseRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Required, validation.Length(3, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   3,
				"max":   50,
				"field": "title",
			},
		})),
		//validation.Field(&reqBody.Address, validation.Required, validation.Length(3, 100).ErrorObject(&errors.ApplicationError{
		//	ErrorType:      errors.RangeValidationErr,
		//	TranslationKey: "rangeValidationError",
		//	TranslationParams: map[string]interface{}{
		//		"min":   3,
		//		"max":   100,
		//		"field": "address",
		//	},
		//})),
	)
}

type WarehouseResponseBody struct {
	ID        uint64 `json:"id"`
	OwnerID   uint64 `json:"owner_id"`
	AccountID uint64 `json:"account_id"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	IsUsed    bool   `json:"is_used"`
}

type WarehouseQueryParams struct {
	KeyWord string
}

type RequstParams struct {
	CreatedBy   uint64
	AccountSlug string
	AccountID   uint64
}
