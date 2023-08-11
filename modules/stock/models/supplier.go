package models

import (
	"github.com/google/uuid"
	commonModels "pi-inventory/common/models"
)

type SupplierResponseBody struct {
	ID                        uuid.UUID               `json:"id"`
	Title                     string                  `json:"title"`
	FirstName                 string                  `json:"first_name"`
	LastName                  string                  `json:"last_name"`
	MiddleName                string                  `json:"middle_name"`
	Suffix                    string                  `json:"suffix"`
	DisplayName               string                  `json:"display_name"`
	CompanyName               string                  `json:"company_name"`
	BusinessID                string                  `json:"business_id"`
	Phone                     string                  `json:"phone"`
	PhoneCountry              string                  `json:"phone_country"`
	Mobile                    string                  `json:"mobile"`
	MobileCountry             string                  `json:"mobile_country"`
	Email                     string                  `json:"email"`
	OtherContact              string                  `json:"other_contact"`
	WebsiteAddress            string                  `json:"website_address"`
	SupplierAddress           string                  `json:"supplier_address"`
	Note                      string                  `json:"note"`
	TaxInformation            string                  `json:"tax_information"`
	SupplierAttachment        []Attachment            `json:"attachments"`
	BillingAddressStreet      string                  `json:"billing_address_street"`
	BillingAddressCity        string                  `json:"billing_address_city"`
	BillingAddressState       string                  `json:"billing_address_state"`
	BillingAddressPostalCode  string                  `json:"billing_address_postal_code"`
	BillingAddressCountry     string                  `json:"billing_address_country"`
	ShippingAddressStreet     string                  `json:"shipping_address_street"`
	ShippingAddressCity       string                  `json:"shipping_address_city"`
	ShippingAddressState      string                  `json:"shipping_address_state"`
	ShippingAddressPostalCode string                  `json:"shipping_address_postal_code"`
	ShippingAddressCountry    string                  `json:"shipping_address_country"`
	Status                    string                  `json:"status"`
	BillingRate               string                  `json:"billing_rate"`
	AsOf                      string                  `json:"as_of"`
	DefaultExpenseAccount     string                  `json:"default_expense_account"`
	OpeningBalance            commonModels.FloatMoney `json:"unit_opening_balance"`
	Language                  string                  `json:"language"`
}

type Attachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}
