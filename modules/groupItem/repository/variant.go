package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	groupItemModuleConst "pi-inventory/modules/groupItem/consts"
	"pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/schema"

	"gorm.io/gorm"
)

type VariantRepositoryInterface interface {
	CreateVariant(variant *schema.Variant) (*schema.Variant, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.VariantQueryParams) ([]*schema.Variant, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Variant, error)
	DeleteVariant(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
	UpdateVariant(variant *schema.Variant) (*schema.Variant, error)
}

type variantRepository struct {
	Db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *variantRepository {
	return &variantRepository{Db: db}
}

func (vr *variantRepository) CreateVariant(variant *schema.Variant) (*schema.Variant, error) {
	if err := vr.Db.Table(groupItemModuleConst.VariantTable).Create(&variant).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return variant, nil
}

func (vr *variantRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.VariantQueryParams) ([]*schema.Variant, error) {
	variants := []*schema.Variant{}
	query := fmt.Sprintf("%s=?", groupItemModuleConst.VariantFieldAccountID)
	baseQuery := vr.Db.Table(groupItemModuleConst.VariantTable).Where(query, requestParams.AccountID).Order("created_at DESC")
	if len(queryParamStruct.KeyWord) != 0 {
		if err := baseQuery.Where("title LIKE ?", "%"+queryParamStruct.KeyWord+"%").Find(&variants).Error; err != nil {
			logger.LogError(err)
			return nil, err
		}
	}
	if err := baseQuery.Find(&variants).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return variants, nil
}

func (vr *variantRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Variant, error) {
	var variant schema.Variant
	query := fmt.Sprintf("%s=? AND %s = ?", groupItemModuleConst.VariantFieldAccountID, field)
	if err := vr.Db.Table(groupItemModuleConst.VariantTable).Where(query, requestParams.AccountID, value).First(&variant).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return &variant, nil
}

func (vr *variantRepository) DeleteVariant(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	var variant schema.Variant
	query := fmt.Sprintf("%s=? AND %s=?", groupItemModuleConst.VariantFieldAccountID, groupItemModuleConst.VariantFieldID)
	if err := vr.Db.Table(groupItemModuleConst.VariantTable).Where(query, requestParams.AccountID, ID).Delete(&variant).Unscoped().Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return ID, nil
}

func (vr *variantRepository) UpdateVariant(variant *schema.Variant) (*schema.Variant, error) {
	if err := vr.Db.Table(groupItemModuleConst.VariantTable).Save(&variant).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return variant, nil
}
