package service

import (
	"context"
	"encoding/json"
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	attachmentModel "pi-inventory/modules/attachment/models"
	attachmentService "pi-inventory/modules/attachment/service"
	"pi-inventory/modules/groupItem/consts"
	"pi-inventory/modules/groupItem/models"
	"pi-inventory/modules/groupItem/repository"
	"pi-inventory/modules/groupItem/schema"
	stockConsts "pi-inventory/modules/stock/consts"
	stockModel "pi-inventory/modules/stock/models"
	stockRepository "pi-inventory/modules/stock/repository"
	stockSchema "pi-inventory/modules/stock/schema"
	stockService "pi-inventory/modules/stock/service"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GroupItemServiceInterface interface {
	CreateGroupItem(requestParams commonModel.AccountRequstParams, reqBody *models.RequestGroupItemBody, context context.Context) (*schema.GroupItem, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.GroupItemQueryParams) ([]*models.SingleGroupItemResponse, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, groupItemID uint64) (*models.ResponseGroupItemBody, error)
	UpdateGroupItem(requestParams commonModel.AccountRequstParams, groupItemID uint64, reqBody *models.UpdateGroupItemRequestBody) (*schema.GroupItem, error)
	GroupItemSummary(requestParams commonModel.AccountRequstParams) (int64, error)
}

type groupItemService struct {
	groupItemRepository repository.GroupItemRepositoryInterface
	unitService         stockService.UnitServiceInterface
	variantService      VariantServiceInterface
	stockService        stockService.StockServiceInterface
	attachmentService   attachmentService.AttachmentServiceInterface
	stockRepository     stockRepository.StockRepositoryInterface
}

func NewGroupItemService(GroupItemRepository repository.GroupItemRepositoryInterface,
	unitService stockService.UnitServiceInterface,
	variantService VariantServiceInterface,
	stockService stockService.StockServiceInterface,
	attachmentService attachmentService.AttachmentServiceInterface,
	stockRepository stockRepository.StockRepositoryInterface) *groupItemService {
	return &groupItemService{
		groupItemRepository: GroupItemRepository,
		unitService:         unitService,
		variantService:      variantService,
		stockService:        stockService,
		attachmentService:   attachmentService,
		stockRepository:     stockRepository,
	}
}

func (gis *groupItemService) CreateGroupItem(requestParams commonModel.AccountRequstParams, reqBody *models.RequestGroupItemBody, context context.Context) (*schema.GroupItem, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	variantJSON, err := convertToJSON(reqBody.Variants)
	if err != nil {
		return nil, err
	}

	newGroupItem := &schema.GroupItem{}
	setReqGroupItemFieldValues(newGroupItem, reqBody, requestParams, variantJSON)

	_, err = gis.groupItemRepository.FindBy(requestParams, consts.GroupItemFieldName, reqBody.Name)
	if err == nil {
		logger.LogError("A groupItem with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "group item",
			},
			HttpCode: http.StatusBadRequest,
		}
	}

	groupItem, err := gis.groupItemRepository.CreateGroupItem(newGroupItem)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Group item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	if len(reqBody.GroupLineItems) > 0 {
		for _, groupLineItem := range reqBody.GroupLineItems {
			newStock := &stockModel.AddStockRequestBody{
				Name:           groupLineItem.Title,
				SKU:            groupLineItem.SKU,
				SellingPrice:   groupLineItem.SellingPrice,
				PurchasePrice:  groupLineItem.CostPrice,
				AsOfDate:       groupLineItem.AsOfDate,
				PurchaseDate:   groupLineItem.PurchaseDate,
				ExpiryDate:     groupLineItem.ExpiryDate,
				TrackInventory: groupLineItem.IsStocked,
				StockQty:       groupLineItem.StockQty,
				ReorderQty:     groupLineItem.ReorderQty,
				LocationID:     groupLineItem.LocationID,
				SupplierID:     groupLineItem.SupplierID,
				GroupItemID:    groupItem.ID,
				AccountID:      groupLineItem.AccountID,
			}
			_, err = gis.stockService.CreateStock(requestParams, newStock, context)
			if err != nil {
				return nil, err
			}
		}
	}
	return groupItem, nil
}

func (gis *groupItemService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.GroupItemQueryParams) ([]*models.SingleGroupItemResponse, *commonModel.PageResponse, error) {

	groupItems, recordCount, err := gis.groupItemRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	var countList []*int64
	for _, singleInstance := range groupItems {
		lineItemsCount, err := gis.stockRepository.CountOfStockList(requestParams, stockConsts.StockFieldGroupItemID, singleInstance.ID, nil)
		if err != nil {
			return nil, nil, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "unknownError",
				HttpCode:       http.StatusInternalServerError,
			}
		}
		countList = append(countList, &lineItemsCount)
	}
	responseGroupItems := setGroupItemResponse(groupItems, countList)
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}
	return responseGroupItems, &pageInfo, nil
}

func (gis *groupItemService) FindByID(requestParams commonModel.AccountRequstParams, groupItemID uint64) (*models.ResponseGroupItemBody, error) {
	singleGroupItem, err := gis.groupItemRepository.FindBy(requestParams, consts.GroupItemFieldID, groupItemID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "group item",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	stocks, err := gis.stockRepository.FindBy(requestParams, stockConsts.StockFieldGroupItemID, singleGroupItem[0].ID)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "UnknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	if len(stocks) == 0 {
		logger.LogError("groupItemID NotFound in stocks table")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.NotFoundErr,
			TranslationKey: "idNotFound",
			TranslationParams: map[string]interface{}{
				"field": "stocks",
			},
			HttpCode: http.StatusBadRequest,
		}
	}
	attachments, err := gis.attachmentService.FetchAttachments(singleGroupItem[0].AttachmentKey)
	if err != nil {
		return nil, err
	}
	// setAttachmentsView(responseStocks[0], attachments)
	resGroupItem := setSingleGroupItemResponse(singleGroupItem, attachments)
	resGroupItem.GroupLineItems = setResponseGroupLineItems(stocks)
	return resGroupItem, nil
}

func (gis *groupItemService) UpdateGroupItem(requestParams commonModel.AccountRequstParams, groupItemID uint64, reqBody *models.UpdateGroupItemRequestBody) (*schema.GroupItem, error) {

	prevGroupItem, err := gis.groupItemRepository.FindBy(requestParams, consts.GroupItemFieldID, groupItemID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "group item",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	if len(reqBody.GroupLineItems) > 0 {
		for _, groupLineItem := range reqBody.GroupLineItems {
			newGroupLineItem := &stockModel.UpdateStockRequestBody{
				Name:           groupLineItem.Title,
				SKU:            groupLineItem.SKU,
				SellingPrice:   groupLineItem.SellingPrice,
				PurchasePrice:  groupLineItem.PurchasePrice,
				TrackInventory: groupLineItem.IsStocked,
				StockQty:       groupLineItem.StockQty,
				ReorderQty:     groupLineItem.ReorderQty,
				AsOfDate:       groupLineItem.AsOfDate,
				PurchaseDate:   groupLineItem.PurchaseDate,
				ExpiryDate:     groupLineItem.ExpiryDate,
				SupplierID:     groupLineItem.SupplierID,
				CategoryID:     groupLineItem.CategoryID,
				LocationID:     groupLineItem.LocationID,
				UnitID:         groupLineItem.UnitID,
				AccountID:      requestParams.AccountID,
			}
			_, err := gis.stockService.UpdateStock(requestParams, groupLineItem.ID, newGroupLineItem)
			if err != nil {
				return nil, err
			}
		}
	}
	updatedGroupItem := updateGroupItemFields(reqBody, prevGroupItem, requestParams)
	groupItem, err := gis.groupItemRepository.UpdateGroupItem(updatedGroupItem)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Group item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	return groupItem, nil
}

func (gis *groupItemService) GroupItemSummary(requestParams commonModel.AccountRequstParams) (int64, error) {
	groupItemCount, err := gis.groupItemRepository.GroupItemSummary(requestParams)
	return groupItemCount, err
}

func setGroupItemResponse(list []*schema.GroupItem, countList []*int64) []*models.SingleGroupItemResponse {

	groupItems := make([]*models.SingleGroupItemResponse, 0)
	for i, singleInstance := range list {
		groupItem := models.SingleGroupItemResponse{}
		setResponseGroupItemFieldValues(singleInstance, &groupItem, countList[i])
		groupItems = append(groupItems, &groupItem)

	}

	return groupItems
}

func setResponseGroupItemFieldValues(instance *schema.GroupItem, resGroupItem *models.SingleGroupItemResponse, lineItemsCount *int64) {
	resGroupItem.ID = instance.ID
	resGroupItem.Name = instance.Name
	resGroupItem.Tag = instance.Tag
	resGroupItem.AccountID = instance.AccountID
	resGroupItem.PurchasePriceCost = commonModel.Money{
		Amount:   instance.PurchasePrice,
		Currency: instance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resGroupItem.SellingPriceCost = commonModel.Money{
		Amount:   instance.SellingPrice,
		Currency: instance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	resGroupItem.CountOfGroupLineItems = int(*lineItemsCount)
}

func convertToJSON(variants []*models.RequestVariant) (datatypes.JSON, error) {
	variantJSON, err := json.Marshal(variants)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(variantJSON), nil
}
func setReqGroupItemFieldValues(newInstance *schema.GroupItem, reqGroupItem *models.RequestGroupItemBody, requestParams commonModel.AccountRequstParams, variantJSON datatypes.JSON) {
	newInstance.OwnerID = requestParams.CreatedBy
	newInstance.AccountID = requestParams.AccountID
	newInstance.Name = reqGroupItem.Name
	newInstance.Tag = reqGroupItem.Tag
	newInstance.GroupItemUnit = reqGroupItem.GroupItemUnit
	newInstance.SellingPrice = reqGroupItem.SellingPriceCost.Amount
	newInstance.SellingPriceCurrency = reqGroupItem.SellingPriceCost.Currency
	newInstance.PurchasePrice = reqGroupItem.PurchasePriceCost.Amount
	newInstance.PurchasePriceCurrency = reqGroupItem.PurchasePriceCost.Currency
	newInstance.AttachmentKey = reqGroupItem.AttachmentKey
	newInstance.Variant_values = variantJSON
	newInstance.CreatedAt = time.Now()
	newInstance.CreatedBy = requestParams.CreatedBy
}

func setSingleGroupItemResponse(newInstance []*schema.GroupItem, attachments []*attachmentModel.AttachmentCustomResponse) *models.ResponseGroupItemBody {
	resGroupItem := &models.ResponseGroupItemBody{}
	resGroupItem.ID = newInstance[0].ID
	resGroupItem.Name = newInstance[0].Name
	resGroupItem.Tag = newInstance[0].Tag
	resGroupItem.AccountID = newInstance[0].AccountID
	resGroupItem.GroupItemUnit = newInstance[0].GroupItemUnit
	resGroupItem.SellingPriceCost = commonModel.Money{
		Amount:   newInstance[0].SellingPrice,
		Currency: newInstance[0].SellingPriceCurrency,
	}.ConvertFloatMoney()
	resGroupItem.PurchasePriceCost = commonModel.Money{
		Amount:   newInstance[0].PurchasePrice,
		Currency: newInstance[0].PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resGroupItem.Variants = newInstance[0].Variant_values
	resGroupItem.Attachments = append(resGroupItem.Attachments, attachments...)
	return resGroupItem
}

func setResponseGroupLineItems(groupLineItems []*stockSchema.Stock) []*models.ResponseGroupLineItem {
	responseGroupLineItems := make([]*models.ResponseGroupLineItem, len(groupLineItems))

	for i, groupLineItem := range groupLineItems {
		responseGroupLineItems[i] = &models.ResponseGroupLineItem{
			Title:      groupLineItem.Name,
			SKU:        groupLineItem.SKU,
			StockQty:   groupLineItem.StockQty,
			ReorderQty: groupLineItem.ReorderQty,
			StockID:    groupLineItem.ID,
			SupplierID: groupLineItem.SupplierID,
			LocationID: groupLineItem.LocationID,
			AccountID:  groupLineItem.AccountID,
			CostPrice: commonModel.Money{
				Amount:   groupLineItem.PurchasePrice,
				Currency: groupLineItem.PurchasePriceCurrency,
			}.ConvertFloatMoney(),
			SellingPrice: commonModel.Money{
				Amount:   groupLineItem.SellingPrice,
				Currency: groupLineItem.SellingPriceCurrency,
			}.ConvertFloatMoney(),
			AsOfDate:     groupLineItem.AsOfDate,
			PurchaseDate: groupLineItem.PurchaseDate,
			ExpiryDate:   groupLineItem.ExpiryDate,
			IsStocked:    groupLineItem.TrackInventory,
		}
	}

	return responseGroupLineItems
}

func updateGroupItemFields(reqBody *models.UpdateGroupItemRequestBody, prevGroupItem []*schema.GroupItem, requestParams commonModel.AccountRequstParams) *schema.GroupItem {
	if reqBody.Name != "" {
		prevGroupItem[0].Name = reqBody.Name
	}
	if reqBody.Tag != "" {
		prevGroupItem[0].Tag = reqBody.Tag
	}
	if reqBody.GroupItemUnit != "" {
		prevGroupItem[0].GroupItemUnit = reqBody.GroupItemUnit
	}
	if reqBody.SellingPriceCost.Amount.GreaterThan(decimal.NewFromFloat32(0)) {
		prevGroupItem[0].SellingPrice = reqBody.SellingPriceCost.Amount
	}
	if reqBody.PurchasePriceCost.Amount.GreaterThan(decimal.NewFromFloat32(0)) {
		prevGroupItem[0].PurchasePrice = reqBody.PurchasePriceCost.Amount
	}
	if reqBody.SellingPriceCost.Currency != "" {
		prevGroupItem[0].SellingPriceCurrency = reqBody.SellingPriceCost.Currency
	}
	if reqBody.PurchasePriceCost.Currency != "" {
		prevGroupItem[0].PurchasePriceCurrency = reqBody.PurchasePriceCost.Currency
	}
	if len(reqBody.AttachmentKey) > 0 {
		prevGroupItem[0].AttachmentKey = reqBody.AttachmentKey
	}
	prevGroupItem[0].UpdatedAt = time.Now()
	prevGroupItem[0].UpdatedBy = requestParams.CreatedBy
	return prevGroupItem[0]
}
