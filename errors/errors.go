package errors

import (
	"errors"
	"fmt"
	"net/http"

	"pi-inventory/common/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorType int
type Body map[string]interface{}

const (
	NotFoundErr ErrorType = iota + 1
	InvalidEmail
	UnKnownErr
	InvalidRequestParams
	RecordNotFound
	Unsuccessfull
	RangeValidationErr
	AlreadyExistsErr
	// InvalidDate
	// Overflow
	InvalidType
	// InvalidOrder
	// MisMatchSubtotal
	// MisMatchTotal
	// ZeroEstimateLineItems
	// InvalidLineItemType
	RequiredField
	// InvalidAmount
	LengthValidation
	// Underflow
	// InvaidPhone
	// InvaidCountryCode
	// InvaidMobile
	// InvalidStatus
	// InvalidOrderBy
	// InvalidOrderType
	// ExcelFieldError
	// ExcelFieldMissing
	// ExcelMissingData
)

type ApplicationError struct {
	ErrorType         ErrorType
	TranslationKey    string
	TranslationParams map[string]interface{}
	HttpCode          int
	Errs              []ApplicationError
}

func (e *ApplicationError) Join(errs ...ApplicationError) {
	e.Errs = append(e.Errs, errs...)
	//return e
}

func (e *ApplicationError) Unwrap() []ApplicationError {
	return e.Errs
}

func (e *ApplicationError) Error() string {
	return utils.Trans(e.TranslationKey, e.TranslationParams)
}

func (e *ApplicationError) Code() string {
	return e.TranslationKey
}

func (e *ApplicationError) Message() string {
	return e.TranslationKey
}

func (e *ApplicationError) SetMessage(message string) validation.Error {
	e.TranslationKey = message
	return e
}

func (e *ApplicationError) Params() map[string]interface{} {
	return e.TranslationParams
}

func (e *ApplicationError) SetParams(params map[string]interface{}) validation.Error {
	e.TranslationParams = params
	return e
}

type GinError struct{}

func GetErrorMessage(err error) string {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {
		return utils.Trans(aError.TranslationKey, aError.TranslationParams)
	}
	return fmt.Sprint(err)
}

func (g *GinError) GetErrorMessage(err error) map[string]interface{} {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {

		return GenerateErrorResponseBody(err)
	}
	return ValidationErrors(err)
}

func (g *GinError) GetStatusCode(err error) int {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {
		return aError.HttpCode
	}
	return http.StatusBadRequest
}

func GetStatusCode(err error) int {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {
		return aError.HttpCode
	}
	return 500
}

func GetErrorType(err error) int {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {
		return int(aError.ErrorType)
	}
	return 1000
}

func (g *GinError) ErrorTraverse(err error) Body {
	aError := &ApplicationError{}
	if errors.As(err, &aError) {
		if len(aError.Errs) == 0 {
			return GenerateErrorResponseBody(err)
		}

		mp := make(map[string]interface{}, 0)
		for _, e := range aError.Errs {
			field, _ := e.TranslationParams["field"].(string)
			if len(e.Errs) > 0 {
				childErr := g.ErrorTraverse(&e)
				for k, ee := range childErr {
					mp[k] = ee
				}
			} else {
				mp[field] = e.Error()
			}
		}

		field, _ := aError.TranslationParams["field"].(string)
		errBody := map[string]interface{}{field: mp}
		return errBody
	} else {
		return validationResponse(err)
	}
}

func GenerateErrorResponseBody(err error) Body {
	message := err.Error()
	return readFromMap(message)
}

func validationResponse(err error) Body {
	return Body{
		"validation_error": err,
	}
}

func readFromMap(message string) Body {
	return GenerateResponseBody(message)
}

func GenerateResponseBody(message string) Body {
	return Body{
		"message": message,
	}
}

func ValidationErrors(err error) Body {
	return validationResponse(err)
}
