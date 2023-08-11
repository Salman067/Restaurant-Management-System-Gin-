package repository

import "gorm.io/gorm"

type SupplierRepositoryInterface interface {
}

type supplierRepository struct {
	Db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *supplierRepository {
	return &supplierRepository{Db: db}
}
