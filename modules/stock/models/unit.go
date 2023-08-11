package models

import (
	"pi-inventory/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateUnitRequestBody struct {
	Title     string `json:"title"`
	Status    string `json:"status"`
	IsUsed    bool   `json:"is_used"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody UpdateUnitRequestBody) Validate() error {
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

type AddUnitRequestBody struct {
	Title     string `json:"title"`
	IsUsed    bool   `json:"is_used"`
	AccountID uint64 `json:"account_id"`
}

func (reqBody AddUnitRequestBody) Validate() error {
	return validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Title, validation.Required, validation.Length(2, 50).ErrorObject(&errors.ApplicationError{
			ErrorType:      errors.RangeValidationErr,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"min":   2,
				"max":   50,
				"field": "title",
			},
		})),
	)
}

type UnitResponseBody struct {
	ID        uint64 `json:"id"`
	OwnerID   uint64 `json:"owner_id"`
	AccountID uint64 `json:"account_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	IsUsed    bool   `json:"is_used"`
}

type UnitQueryParams struct {
	KeyWord string
}
