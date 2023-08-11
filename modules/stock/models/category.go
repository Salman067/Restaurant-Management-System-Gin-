package models

import (
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateCategoryRequestBody struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Description string `json:"description"`
	IsUsed      bool   `json:"is_used"`
	AccountID   uint64 `json:"account_id"`
}

func (reqBody UpdateCategoryRequestBody) Validate() error {
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

type AddCategoryRequestBody struct {
	Title       string `json:"title"`
	AccountID   uint64 `json:"account_id"`
	Description string `json:"description"`
	IsUsed      bool   `json:"is_used"`
}

func (reqBody AddCategoryRequestBody) Validate() error {
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
	)
}

type CategoryResponseBody struct {
	ID          uint64 `json:"id"`
	OwnerID     uint64 `json:"owner_id"`
	AccountID   uint64 `json:"account_id"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	Description string `json:"description"`
	IsUsed      bool   `json:"is_used"`
}

type CategoryQueryParams struct {
	KeyWord string
}
