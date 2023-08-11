package repository

import (
	"gorm.io/gorm"
)

type ProfileRepositoryInterface interface {
	//FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.WarehouseQueryParams) ([]*schema.Warehouse, int64, error)
	//Create(warehouse *schema.Warehouse) (*schema.Warehouse, error)
	//Update(warehouse *schema.Warehouse) (*schema.Warehouse, error)
	//Delete(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	//FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Warehouse, error)
}

type profileRepository struct {
	Db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepository {
	return &profileRepository{Db: db}
}
