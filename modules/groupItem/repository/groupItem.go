package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/modules/groupItem/consts"
	groupItemModel "pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/schema"
)

type GroupItemRepositoryInterface interface {
	CreateGroupItem(groupItem *schema.GroupItem) (*schema.GroupItem, error)
	UpdateGroupItem(groupItem *schema.GroupItem) (*schema.GroupItem, error)
	FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.GroupItem, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *groupItemModel.GroupItemQueryParams) ([]*schema.GroupItem, int64, error)
	GroupItemSummary(requestParams commonModel.AccountRequstParams) (int64, error)
}

type groupItemRepository struct {
	Db *gorm.DB
}

func NewGroupItemRepository(db *gorm.DB) *groupItemRepository {
	return &groupItemRepository{Db: db}
}

func (gir *groupItemRepository) CreateGroupItem(groupItem *schema.GroupItem) (*schema.GroupItem, error) {
	if err := gir.Db.Clauses(clause.Returning{}).Table(consts.GroupItemTable).Create(&groupItem).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return groupItem, nil
}

func (gir *groupItemRepository) UpdateGroupItem(groupItem *schema.GroupItem) (*schema.GroupItem, error) {
	err := gir.Db.Clauses(clause.Returning{}).Table(consts.GroupItemTable).Save(&groupItem).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return groupItem, nil
}

func (gir *groupItemRepository) FindBy(requestParams commonModel.AccountRequstParams, field string, value any) ([]*schema.GroupItem, error) {
	groupItem := make([]*schema.GroupItem, 0)
	query := fmt.Sprintf("%s=? AND %s = ? ", consts.GroupItemFieldAccountID, field)
	if err := gir.Db.Model(&schema.GroupItem{}).
		Where(query, requestParams.AccountID, value).
		First(&groupItem).Error; err != nil {
		logger.LogError(err)
		return nil, err
	}
	return groupItem, nil
}

func (gir *groupItemRepository) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *groupItemModel.GroupItemQueryParams) ([]*schema.GroupItem, int64, error) {
	var count int64
	groupItems := []*schema.GroupItem{}
	query1 := fmt.Sprintf("%s=?", consts.GroupItemFieldAccountID)
	if err := gir.Db.Table(consts.GroupItemTable).Where(query1, requestParams.AccountID).Count(&count).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}

	baseQuery := gir.Db.Model(&schema.GroupItem{}).
		Offset(requestParams.Page.Offset).
		Limit(requestParams.Page.Limit).
		Where(query1, requestParams.AccountID).
		Order("created_at DESC")

	// Full text search
	//query := fmt.Sprintf("to_tsvector(%s || ' '|| %s || ' '|| %s || ' '|| %s || ' '|| %s) @@ to_tsquery(?)",
	//	consts.GroupItemFieldName, consts.GroupItemFieldGroupItemUnit, consts.GroupItemFieldTag, consts.GroupItemFieldSellingPriceCost, consts.GroupItemFieldPurchasePriceCost)
	//if len(queryParamStruct.KeyWord) != 0 {
	//	if err := baseQuery.Where(query, queryParamStruct.KeyWord).Find(&groupItems).Error; err != nil {
	//		logger.LogError(err)
	//		return nil, 0, err
	//	}
	//}
	// Partial search on Name, GIUnit, Tag, Selling Cost and Purchase Cost
	query := fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?",
		consts.GroupItemFieldName, consts.GroupItemFieldGroupItemUnit, consts.GroupItemFieldTag, consts.GroupItemFieldSellingPriceCost, consts.GroupItemFieldPurchasePriceCost)

	if len(queryParamStruct.KeyWord) != 0 {
		keyword := fmt.Sprintf("%%%s%%", queryParamStruct.KeyWord) // Add wildcards to the keyword
		if err := baseQuery.Where(query, keyword, keyword, keyword, keyword, keyword).Find(&groupItems).Count(&count).Error; err != nil {
			logger.LogError(err)
			return nil, 0, err
		}
	}

	if err := baseQuery.Find(&groupItems).Error; err != nil {
		logger.LogError(err)
		return nil, 0, err
	}
	return groupItems, count, nil
}

func (gir *groupItemRepository) GroupItemSummary(requestParams commonModel.AccountRequstParams) (int64, error) {
	var count int64
	query := fmt.Sprintf("%s=?", consts.GroupItemFieldAccountID)
	if err := gir.Db.Table(consts.GroupItemTable).Count(&count).Where(query, requestParams.AccountID).Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return count, nil
}
