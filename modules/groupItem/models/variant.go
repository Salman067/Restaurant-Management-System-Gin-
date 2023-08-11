package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"pi-inventory/errors"
)

type UpdateVariantRequestBody struct {
	Title     string `json:"title"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody UpdateVariantRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Length(3, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"field": "title",
				"min":   3,
				"max":   50,
			},
		})),
	)
}

type AddVariantRequestBody struct {
	Title     string `json:"title"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody AddVariantRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Required, validation.Length(3, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"field": "title",
				"min":   3,
				"max":   50,
			},
		})),
	)
}

type VariantResponseBody struct {
	ID        uint64 `json:"id"`
	OwnerID   uint64 `json:"owner_id"`
	AccountID uint64 `json:"account_id"`
	Title     string `json:"title"`
}

type VariantQueryParams struct {
	KeyWord string
}
