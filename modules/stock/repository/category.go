package repository

import (
	"fmt"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	stockConst "pi-inventory/modules/stock/consts"
	stockModel "pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/schema"

	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.CategoryQueryParams) (*[]schema.Category, int64, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Category, error)
	Create(category *schema.Category) (*schema.Category, error)
	Update(category *schema.Category) (*schema.Category, error)
	Delete(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
}

type categoryRepository struct {
	Db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{Db: db}
}

func (cr *categoryRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.CategoryQueryParams) (*[]schema.Category, int64, error) {
	var count int64
	category := []schema.Category{}
	query := fmt.Sprintf("%s=?", stockConst.CategoryFieldAccountID)
	if err := cr.Db.Table(stockConst.CategoryTable).Where(query, requestParams.AccountID).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	baseQuery := cr.Db.Table(stockConst.CategoryTable).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Where(query, requestParams.AccountID).
		Order("created_at DESC")

	// full text search  on Title and description field
	//query := fmt.Sprintf("to_tsvector(%s || ' '|| %s) @@ to_tsquery(?)",
	//	stockConst.CategoryFieldTitle, stockConst.CategoryFieldDescription)
	//
	//if len(queryParamStruct.KeyWord) != 0 {
	//	if err := baseQuery.Where(query, queryParamStruct.KeyWord).Find(&category).Error; err != nil {
	//		logger.LogError(err)
	//		return nil, 0, err
	//	}
	//}

	// Partial search on Title and description field
	searchQ := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
		stockConst.CategoryFieldTitle, stockConst.CategoryFieldTitle, stockConst.CategoryFieldDescription)
	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		if err := baseQuery.Where(searchQ, keyword, keyword, keyword).Find(&category).Count(&count).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
	}

	if err := baseQuery.Find(&category).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return &category, count, nil
}

func (cr *categoryRepository) Create(category *schema.Category) (*schema.Category, error) {
	err := cr.Db.Table(stockConst.CategoryTable).Create(&category).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) Update(category *schema.Category) (*schema.Category, error) {
	err := cr.Db.Table(stockConst.CategoryTable).Save(&category).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) Delete(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {
	category := schema.Category{}
	query := fmt.Sprintf("%s=? AND %s = ? ", stockConst.CategoryFieldAccountID, stockConst.CategoryFieldID)
	err := cr.Db.Table(stockConst.CategoryTable).Where(query, requestParams.AccountID, ID).Delete(&category).Error
	if err != nil {
		logger.LogError(err)
		return 0, err
	}
	return ID, nil
}

func (cr *categoryRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) (*schema.Category, error) {
	category := &schema.Category{}
	query := fmt.Sprintf("%s=? AND %s = ? ", stockConst.CategoryFieldAccountID, field)
	if err := cr.Db.Table(stockConst.CategoryTable).Where(query, requestParams.AccountID, value).First(&category).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return category, nil
}
