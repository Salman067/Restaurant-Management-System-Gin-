package models

import (
	"math"
	"pi-inventory/errors"
	"strconv"

	"github.com/shopspring/decimal"
)

type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

type FloatMoney struct {
	Amount      decimal.Decimal `json:"-"`
	FloatAmount float64         `json:"amount"`
	Currency    string          `json:"currency"`
}

func (m Money) ConvertFloatMoney() FloatMoney {
	amount := 0.0
	amount, _ = strconv.ParseFloat(m.Amount.String(), 64)
	return FloatMoney{
		Amount:      m.Amount,
		FloatAmount: amount,
		Currency:    m.Currency,
	}
}

func ValidateMoney(value interface{}) error {
	money, ok := value.(Money)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "money",
			},
		}
	}
	if money.Amount.LessThan(decimal.Zero) {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "lessThanError",
			TranslationParams: map[string]interface{}{
				"field": "amount",
				"value": 0,
			},
		}
	}
	if len(money.Currency) == 0 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"field": "currency",
				"min":   2,
				"max":   12,
			},
		}
	}
	return nil
}

func ValidateMoneyAmount(money Money, errorField string) error {

	moneyAmountString := money.Amount.Abs().String()
	numberOfDigitsInMoneyAmountString := len(moneyAmountString)
	if numberOfDigitsInMoneyAmountString > 10 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "moneyOverflowed",
			TranslationParams: map[string]interface{}{
				"value": 11,
			},
		}
	}
	fractionalDigitsOfMoneyAmount := math.Abs(float64(money.Amount.Exponent()))
	if fractionalDigitsOfMoneyAmount > 2 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "moneyOverflowed1",
			TranslationParams: map[string]interface{}{
				"value": 3,
			},
		}
	}
	return nil
}

func ValidateNewSellingPrice(value interface{}) error {
	money, ok := value.(Money)
	if !ok {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "invalidValueType",
			TranslationParams: map[string]interface{}{
				"field": "money",
			},
		}
	}
	if money.Amount.LessThanOrEqual(decimal.Zero) {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "greaterThanError",
			TranslationParams: map[string]interface{}{
				"field": "amount",
				"value": 0,
			},
		}
	}
	if len(money.Currency) == 0 {
		return &errors.ApplicationError{
			ErrorType:      errors.RequiredField,
			TranslationKey: "rangeValidationError",
			TranslationParams: map[string]interface{}{
				"field": "currency",
				"min":   2,
				"max":   12,
			},
		}
	}
	return nil
}
