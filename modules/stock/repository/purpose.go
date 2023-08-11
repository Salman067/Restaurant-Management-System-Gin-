package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	stockConst "pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/schema"

	"gorm.io/gorm"
)

type PurposeRepositoryInterface interface {
	CreatePurpose(purpose *schema.Purpose) (*schema.Purpose, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.PurposeQueryParams) ([]*schema.Purpose, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Purpose, error)
	DeletePurpose(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	UpdatePurpose(purpose *schema.Purpose) (*schema.Purpose, error)
}

type purposeRepository struct {
	Db *gorm.DB
}

func NewPurposeRepository(db *gorm.DB) *purposeRepository {
	return &purposeRepository{Db: db}
}

func (pr *purposeRepository) CreatePurpose(purpose *schema.Purpose) (*schema.Purpose, error) {
	if err := pr.Db.Table(stockConst.PurposeTable).Create(&purpose).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return purpose, nil
}

func (pr *purposeRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.PurposeQueryParams) ([]*schema.Purpose, error) {
	purposes := []*schema.Purpose{}
	query := fmt.Sprintf("%s=?", stockConst.PurposeFieldAccountID)
	baseQuery := pr.Db.Table(stockConst.PurposeTable).Where(query, requestParams.AccountID).Order("created_at DESC")

	// Partial Search
	searchQ := fmt.Sprintf("%s ILIKE ?", stockConst.PurposeFieldTitle)
	if len(queryParamStruct.KeyWord) != 0 {
		if err := baseQuery.Where(searchQ, "%"+queryParamStruct.KeyWord+"%").Find(&purposes).Error; err != nil {
			logger.LogError(err)
			return nil, err
		}
	}
	if err := baseQuery.Find(&purposes).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return purposes, nil
}

func (pr *purposeRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Purpose, error) {
	var purpose schema.Purpose
	query := fmt.Sprintf("%s=? AND %s = ?", stockConst.PurposeFieldAccountID, field)
	if err := pr.Db.Table(stockConst.PurposeTable).Where(query, requestParams.AccountID, value).First(&purpose).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return &purpose, nil
}

func (pr *purposeRepository) DeletePurpose(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	var purpose schema.Purpose
	query := fmt.Sprintf("%s=? AND %s=?", stockConst.PurposeFieldAccountID, stockConst.PurposeFieldID)
	if err := pr.Db.Table(stockConst.PurposeTable).Where(query, requestParams.AccountID, ID).Delete(&purpose).Unscoped().Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return ID, nil
}

func (pr *purposeRepository) UpdatePurpose(purpose *schema.Purpose) (*schema.Purpose, error) {
	if err := pr.Db.Table(stockConst.PurposeTable).Save(&purpose).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return purpose, nil
}
