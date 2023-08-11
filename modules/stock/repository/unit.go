package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	stockModuleConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/schema"

	"gorm.io/gorm"
)

type UnitRepositoryInterface interface {
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Unit, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.UnitQueryParams) ([]*schema.Unit, int64, error)
	CreateUnit(unit *schema.Unit) (*schema.Unit, error)
	DeleteUnit(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	UpdateUnit(unit *schema.Unit) (*schema.Unit, error)
}

type unitRepository struct {
	Db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *unitRepository {
	return &unitRepository{Db: db}
}

func (ur *unitRepository) CreateUnit(unit *schema.Unit) (*schema.Unit, error) {
	if err := ur.Db.Create(&unit).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return unit, nil
}

func (ur *unitRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.UnitQueryParams) ([]*schema.Unit, int64, error) {
	var units []*schema.Unit
	var count int64
	query := fmt.Sprintf("%s=?", stockModuleConst.UnitFieldAccountID)
	if err := ur.Db.Table(stockModuleConst.UnitTable).Where(query, requestParams.AccountID).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	baseQuery := ur.Db.Table(stockModuleConst.UnitTable).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Where(query, requestParams.AccountID).
		Order("created_at DESC")

	if len(queryParamStruct.KeyWord) != 0 {
		if err := baseQuery.Where("title LIKE ?", "%"+queryParamStruct.KeyWord+"%").Find(&units).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
	}
	if err := baseQuery.Find(&units).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return units, count, nil
}

func (ur *unitRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Unit, error) {
	var unit schema.Unit
	query := fmt.Sprintf("%s = ? AND %s = ?", stockModuleConst.UnitFieldAccountID, field)
	if err := ur.Db.Table(stockModuleConst.UnitTable).Where(query, requestParams.AccountID, value).First(&unit).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return &unit, nil
}

func (ur *unitRepository) DeleteUnit(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	var unit schema.Unit
	query := fmt.Sprintf("%s = ? AND %s = ?", stockModuleConst.UnitFieldAccountID, stockModuleConst.UnitFieldID)
	if err := ur.Db.Table(stockModuleConst.UnitTable).Where(query, requestParams.AccountID, ID).Delete(&unit).Unscoped().Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return ID, nil
}

func (ur *unitRepository) UpdateUnit(unit *schema.Unit) (*schema.Unit, error) {
	if err := ur.Db.Table(stockModuleConst.UnitTable).Save(&unit).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return unit, nil
}
