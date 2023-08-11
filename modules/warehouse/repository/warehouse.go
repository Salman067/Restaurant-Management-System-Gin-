package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	warehouseConst "pi-inventory/modules/warehouse/consts"
	"pi-inventory/modules/warehouse/models"
	"pi-inventory/modules/warehouse/schema"

	"gorm.io/gorm"
	commonModel "pi-inventory/common/models"
)

type WarehouseRepositoryInterface interface {
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.WarehouseQueryParams) ([]*schema.Warehouse, int64, error)
	Create(warehouse *schema.Warehouse) (*schema.Warehouse, error)
	Update(warehouse *schema.Warehouse) (*schema.Warehouse, error)
	Delete(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Warehouse, error)
}

type warehouseRepository struct {
	Db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) *warehouseRepository {
	return &warehouseRepository{Db: db}
}

func (wr *warehouseRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.WarehouseQueryParams) ([]*schema.Warehouse, int64, error) {
	warehouse := []*schema.Warehouse{}
	var count int64
	query := fmt.Sprintf("%s=?", warehouseConst.WarehouseFieldAccountID)
	if err := wr.Db.Table(warehouseConst.WarehouseTable).Where(query, requestParams.AccountID).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	baseQuery := wr.Db.Table(warehouseConst.WarehouseTable).
		Where(query, requestParams.AccountID).
		Order("created_at DESC").
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit)

	searchQ := fmt.Sprintf("%s ILIKE ?", warehouseConst.WarehouseFieldTitle)
	if len(queryParamStruct.KeyWord) != 0 {
		if err := baseQuery.Where(searchQ, "%"+queryParamStruct.KeyWord+"%").Find(&warehouse).Count(&count).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
	}

	if err := baseQuery.Find(&warehouse).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return warehouse, count, nil
}

func (wr *warehouseRepository) Create(warehouse *schema.Warehouse) (*schema.Warehouse, error) {
	logger.LogError(warehouse.AccountID)
	err := wr.Db.Table(warehouseConst.WarehouseTable).Create(&warehouse).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return warehouse, nil
}

func (wr *warehouseRepository) Update(warehouse *schema.Warehouse) (*schema.Warehouse, error) {
	err := wr.Db.Table(warehouseConst.WarehouseTable).Save(&warehouse).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return warehouse, nil
}

func (wr *warehouseRepository) Delete(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	warehouse := []schema.Warehouse{}
	query := fmt.Sprintf("%s=? AND %s=?", warehouseConst.WarehouseFieldID, warehouseConst.WarehouseFieldAccountID)
	if err := wr.Db.Table(warehouseConst.WarehouseTable).Where(query, ID, requestParams.AccountID).Delete(&warehouse).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return ID, nil
}

func (wr *warehouseRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Warehouse, error) {
	warehouse := &schema.Warehouse{}
	query := fmt.Sprintf("%s = ? ", field)
	query1 := fmt.Sprintf("%s=?", warehouseConst.WarehouseFieldAccountID)
	if err := wr.Db.Table(warehouseConst.WarehouseTable).Where(query1, requestParams.AccountID).Where(query, value).First(&warehouse).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return warehouse, nil
}
