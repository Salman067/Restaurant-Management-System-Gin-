package models

import "github.com/google/uuid"

type TaxResponseBody struct {
	ID                  uuid.UUID          `json:"id"`
	TaxName             string             `json:"tax_name"`
	AgencyID            uuid.UUID          `json:"agency_id"`
	Description         string             `json:"description"`
	SalesRateHistory    []*RespRateHistory `json:"sales_tax_rate_history"`
	PurchaseRateHistory []*RespRateHistory `json:"purchase_tax_rate_history"`
	Status              string             `json:"status"`
	CreatedBy           int                `json:"created_by"`
	AccountID           int                `json:"account_id"`
	ParentIds           []uuid.UUID        `json:"parent_ids,omitempty"`
	ChildrensIds        []*ChildrenID      `json:"childrens_ids,omitempty"`
}

type RespRateHistory struct {
	Rate      float64 `json:"rate,omitempty"`
	StartDate string  `json:"start_date,omitempty"`
}

type ChildrenID struct {
	ID      uuid.UUID `json:"id,omitempty"`
	TaxType string    `json:"tax_type,omitempty"`
}
