package repository

import "gorm.io/gorm"

type TaxRepositoryInterface interface {
}

type TaxRepository struct {
	Db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) *TaxRepository {
	return &TaxRepository{Db: db}
}
