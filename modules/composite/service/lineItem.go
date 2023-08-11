package service

import (
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/common/utils"
	"pi-inventory/errors"
	lineItemConst "pi-inventory/modules/composite/consts"
	"pi-inventory/modules/composite/models"
	"pi-inventory/modules/composite/repository"
	"pi-inventory/modules/composite/schema"
	"pi-inventory/modules/stock/service"
	"time"
)

type LineItemServiceInterface interface {
	CreateLineItem(requestParams commonModel.AccountRequstParams, reqLineItems *[]models.AddLineItemRequestBody) (string, error)
	FetchLineItems(requestParams commonModel.AccountRequstParams, lineItemKey string) ([]*models.LineItemResponseBody, error)
	DeleteLineItem(requestParams commonModel.AccountRequstParams, lineItemKey string) error
}

type lineItemService struct {
	lineItemRepository repository.LineItemRepositoryInterface
	stockService       service.StockServiceInterface
}

func NewLineItemService(lineItemRepo repository.LineItemRepositoryInterface, stockService service.StockServiceInterface) *lineItemService {
	return &lineItemService{lineItemRepository: lineItemRepo, stockService: stockService}
}
func (ls *lineItemService) FetchLineItems(requestParams commonModel.AccountRequstParams, lineItemKey string) ([]*models.LineItemResponseBody, error) {
	lineItems, err := ls.lineItemRepository.FindBy(requestParams, lineItemConst.LineItemFieldLineItemKey, lineItemKey)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	reslineItems, err := setLineItemResponse(lineItems)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	return reslineItems, nil
}

func (ls *lineItemService) CreateLineItem(requestParams commonModel.AccountRequstParams, reqLineItems *[]models.AddLineItemRequestBody) (string, error) {
	lineItemKey, err := utils.GenerateKeyHash()
	if err != nil {
		logger.LogInfo(err)
		return "", &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "ErrorWhileGeneratingHash",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	list := make([]*schema.LineItem, 0)
	listOfReqLineItems := *reqLineItems
	for _, reqLineItem := range listOfReqLineItems {
		newLineItem := &schema.LineItem{}
		setRequestLineItemFieldValues(newLineItem, reqLineItem, requestParams, lineItemKey)
		list = append(list, newLineItem)
	}
	lineItemList := *reqLineItems
	_, err = ls.stockService.FindByID(requestParams, lineItemList[0].StockID)
	if err != nil {
		logger.LogError(err)
		return "", err
	}
	err = ls.lineItemRepository.Create(list)
	if err != nil {
		logger.LogError(err)
		return "", &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Line item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return lineItemKey, nil
}

func (ls *lineItemService) DeleteLineItem(requestParams commonModel.AccountRequstParams, lineItemKey string) error {
	lineItems, err := ls.lineItemRepository.FindBy(requestParams, lineItemConst.LineItemFieldLineItemKey, lineItemKey)
	if err != nil {
		return &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unKnownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for _, lineItem := range lineItems {
		lineItem.IsDeleted = true
	}
	err = ls.lineItemRepository.Delete(lineItems)
	if err != nil {
		logger.LogError(err)
		return &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Line item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func setRequestLineItemFieldValues(newInstance *schema.LineItem, resLineItem models.AddLineItemRequestBody, requestParams commonModel.AccountRequstParams, lineItemKey string) {

	newInstance.Title = resLineItem.Title
	newInstance.OwnerID = requestParams.CreatedBy
	newInstance.AccountID = requestParams.AccountID
	newInstance.StockID = resLineItem.StockID
	newInstance.Quantity = resLineItem.Quantity
	newInstance.UnitRate = resLineItem.UnitRate
	newInstance.TotalSellingPrice = resLineItem.SellingPrice.Amount
	newInstance.SellingPriceCurrency = resLineItem.SellingPrice.Currency
	newInstance.PurchaseRate = resLineItem.PurchaseRate
	newInstance.PurchasePriceCurrency = resLineItem.PurchasePrice.Currency
	newInstance.TotalPurchasePrice = resLineItem.PurchasePrice.Amount
	newInstance.LineItemKey = lineItemKey
	newInstance.CreatedAt = time.Now()
	newInstance.CreatedBy = requestParams.AccountID
}
func setResponseLineItemFieldValues(instance *schema.LineItem, reslineItem *models.LineItemResponseBody) {
	reslineItem.Title = instance.Title
	reslineItem.StockID = instance.StockID
	reslineItem.AccountID = instance.AccountID
	reslineItem.Quantity = instance.Quantity
	reslineItem.UnitRate = instance.UnitRate
	reslineItem.SellingPrice = commonModel.Money{
		Amount:   instance.TotalSellingPrice,
		Currency: instance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	reslineItem.PurchaseRate = instance.PurchaseRate
	reslineItem.PurchasePrice = commonModel.Money{
		Amount:   instance.TotalPurchasePrice,
		Currency: instance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
}
func setLineItemResponse(list []*schema.LineItem) ([]*models.LineItemResponseBody, error) {

	lineItems := make([]*models.LineItemResponseBody, 0)
	for _, singleInstance := range list {
		lineItem := &models.LineItemResponseBody{}
		setResponseLineItemFieldValues(singleInstance, lineItem)
		lineItems = append(lineItems, lineItem)
	}
	return lineItems, nil
}
