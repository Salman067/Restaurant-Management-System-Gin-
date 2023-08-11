package models

import (
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdatePurposeRequestBody struct {
	Title     string `json:"title"`
	IsUsed    bool   `json:"is_used"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody UpdatePurposeRequestBody) Validate() error {
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

type AddPurposeRequestBody struct {
	Title     string `json:"title"`
	IsUsed    bool   `json:"is_used"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody AddPurposeRequestBody) Validate() error {
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

type PurposeResponseBody struct {
	ID        uint64 `json:"id"`
	OwnerID   uint64 `json:"owner_id"`
	AccountID uint64 `json:"account_id"`
	Title     string `json:"title"`
	IsUsed    bool   `json:"is_used"`
}

type PurposeQueryParams struct {
	KeyWord string
}
