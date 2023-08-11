package service

import (
	"net/http"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	"pi-inventory/modules/stock/consts"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
	"pi-inventory/modules/stock/schema"
	"strings"
	"time"
)

type StockActivityServiceInterface interface {
	CreateStockActivity(requestParams commonModel.AccountRequstParams, reqStockActivity *models.AddStockActivityRequestBody) (*schema.StockActivity, error)
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockActivityQueryParams) ([]*models.SingleStockActivityResponse, *commonModel.PageResponse, error)
	CreateBulkStockActivity(requestParams commonModel.AccountRequstParams, reqBulkStockActivity *models.BulkAdjustment) ([]*schema.StockActivity, error)
}
type stockActivityService struct {
	stockActivityRepository repository.StockActivityRepositoryInterface
	stockService            StockServiceInterface
	purposeService          PurposeServiceInterface
}

func NewStockActivityService(stockActivityRepo repository.StockActivityRepositoryInterface,
	stockService StockServiceInterface,
	purposeService PurposeServiceInterface) *stockActivityService {
	return &stockActivityService{
		stockActivityRepository: stockActivityRepo,
		stockService:            stockService,
		purposeService:          purposeService,
	}
}

func (sas *stockActivityService) CreateStockActivity(requestParams commonModel.AccountRequstParams, reqStockActivity *models.AddStockActivityRequestBody) (*schema.StockActivity, error) {
	err := reqStockActivity.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	newStockActivity := getStockActivity(reqStockActivity, requestParams)
	stock, err := sas.stockService.FindByID(requestParams, reqStockActivity.StockID)
	if err != nil {
		return nil, err
	}

	_, err = sas.purposeService.FindByID(requestParams, reqStockActivity.PurposeID)
	if err != nil {
		return nil, err
	}

	if reqStockActivity.AdjustedDate != nil && stock.AsOfDate.After(*reqStockActivity.AdjustedDate) {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "lessThanError",
			TranslationParams: map[string]interface{}{
				"field": "Adjusted date",
				"value": "tracking date(" + stock.AsOfDate.Format("2006-01-02") + ")",
			},
			HttpCode: http.StatusBadRequest,
		}
	}

	if stock.Status == consts.StockInactiveStatus {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "statusInactive",
			TranslationParams: map[string]interface{}{
				"field": "stock",
			},
			HttpCode: http.StatusBadRequest,
		}
	}

	stockActivity, err := sas.stockActivityRepository.CreateStockActivity(newStockActivity)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Stock activity",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	if reqStockActivity.PurposeID != 0 {

		updatePurposeReq := &models.UpdatePurposeRequestBody{
			IsUsed: true,
		}
		_, err = sas.purposeService.UpdatePurpose(requestParams, stockActivity.PurposeID, updatePurposeReq)
		if err != nil {
			return nil, err
		}
	}

	newStock := getUpdatedStock(reqStockActivity, stock)
	_, err = sas.stockService.UpdateStock(requestParams, reqStockActivity.StockID, newStock)
	if err != nil {
		return nil, err
	}

	return stockActivity, nil
}

func getUpdatedStock(activity *models.AddStockActivityRequestBody, stock *models.SingleStockResponseView) *models.UpdateStockRequestBody {
	var updateStockRequest models.UpdateStockRequestBody

	if strings.Compare(activity.Mode, "quantity") == 0 {
		updateStockRequest.StockQty = activity.AdjustedQuantity
		updateStockRequest.SellingPrice.Amount = stock.SellingPrice.Amount
		updateStockRequest.SellingPrice.Currency = stock.SellingPrice.Currency
	}
	if strings.Compare(activity.Mode, "value") == 0 {
		updateStockRequest.SellingPrice = activity.AdjustedSellingPrice
		updateStockRequest.SellingPrice.Currency = activity.AdjustedSellingPrice.Currency
		updateStockRequest.StockQty = stock.StockQty
	}
	updateStockRequest.PurchasePrice.Amount = stock.PurchasePrice.Amount
	updateStockRequest.PurchasePrice.Currency = stock.PurchasePrice.Currency
	updateStockRequest.ReorderQty = stock.ReorderQty
	updateStockRequest.SupplierID = stock.SupplierID
	updateStockRequest.AsOfDate = activity.AdjustedDate
	updateStockRequest.PurchaseDate = activity.PurchaseDate
	updateStockRequest.LocationID = activity.LocationID
	updateStockRequest.TrackInventory = true

	return &updateStockRequest
}

func getStockActivity(reqStockActivity *models.AddStockActivityRequestBody, requestParams commonModel.AccountRequstParams) *schema.StockActivity {
	newStockActivity := &schema.StockActivity{
		OwnerID:                       requestParams.CreatedBy,
		AccountID:                     requestParams.AccountID,
		Mode:                          reqStockActivity.Mode,
		OperationType:                 reqStockActivity.OperationType,
		StockID:                       reqStockActivity.StockID,
		PurposeID:                     reqStockActivity.PurposeID,
		LocationID:                    reqStockActivity.LocationID,
		QuantityOnHand:                reqStockActivity.QuantityOnHand,
		NewQuantity:                   reqStockActivity.NewQuantity,
		AdjustedQuantity:              reqStockActivity.AdjustedQuantity,
		AdjustedDate:                  reqStockActivity.AdjustedDate,
		PurchasePreviousValue:         reqStockActivity.PreviousPurchasePrice.Amount,
		PurchasePreviousValueCurrency: reqStockActivity.PreviousPurchasePrice.Currency,
		PurchaseNewValue:              reqStockActivity.NewPurchasePrice.Amount,
		PurchaseNewValueCurrency:      reqStockActivity.NewPurchasePrice.Currency,
		PurchaseAdjustedValue:         reqStockActivity.AdjustedPurchasePrice.Amount,
		PurchaseAdjustedValueCurrency: reqStockActivity.AdjustedPurchasePrice.Currency,
		SellingPreviousValue:          reqStockActivity.PreviousSellingPrice.Amount,
		SellingPreviousValueCurrency:  reqStockActivity.PreviousSellingPrice.Currency,
		SellingNewValue:               reqStockActivity.NewSellingPrice.Amount,
		SellingNewValueCurrency:       reqStockActivity.NewSellingPrice.Currency,
		SellingAdjustedValue:          reqStockActivity.AdjustedSellingPrice.Amount,
		SellingAdjustedValueCurrency:  reqStockActivity.AdjustedSellingPrice.Currency,
		CreatedAt:                     time.Now(),
		CreatedBy:                     requestParams.CreatedBy,
	}
	return newStockActivity
}

func (sas *stockActivityService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.StockActivityQueryParams) ([]*models.SingleStockActivityResponse, *commonModel.PageResponse, error) {
	stockActivities, recordCount, err := sas.stockActivityRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError}
	}
	responseStockActivities := setStockActivityResponse(stockActivities)
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}

	return responseStockActivities, &pageInfo, nil
}

func setStockActivityResponse(activities []*models.CustomStockActivityResponse) []*models.SingleStockActivityResponse {

	stockActivities := make([]*models.SingleStockActivityResponse, 0)

	for _, singleInstance := range activities {
		stockActivity := models.SingleStockActivityResponse{}
		setResponseStockActivityFieldValues(singleInstance, &stockActivity)
		stockActivities = append(stockActivities, &stockActivity)
	}

	return stockActivities
}

func setResponseStockActivityFieldValues(instance *models.CustomStockActivityResponse, resStockActivity *models.SingleStockActivityResponse) {
	resStockActivity.ID = instance.ID
	resStockActivity.OwnerID = instance.OwnerID
	resStockActivity.AccountID = instance.AccountID
	resStockActivity.Mode = instance.Mode
	resStockActivity.OperationType = instance.OperationType
	resStockActivity.StockID = instance.StockID
	resStockActivity.StockTitle = instance.StockTitle
	resStockActivity.PurposeID = instance.PurposeID
	resStockActivity.PurposeTitle = instance.PurposeTitle
	resStockActivity.NewQuantity = instance.NewQuantity
	resStockActivity.AdjustedQuantity = instance.AdjustedQuantity
	resStockActivity.AdjustedDate = instance.AdjustedDate
	resStockActivity.PurchaseDate = instance.PurchaseDate
	resStockActivity.NewSellingPrice = commonModel.Money{
		Amount:   instance.SellingNewValue,
		Currency: instance.SellingNewValueCurrency,
	}.ConvertFloatMoney()
	resStockActivity.AdjustedSellingPrice = commonModel.Money{
		Amount:   instance.SellingAdjustedValue,
		Currency: instance.SellingAdjustedValueCurrency,
	}.ConvertFloatMoney()

}

func (sas *stockActivityService) CreateBulkStockActivity(requestParams commonModel.AccountRequstParams, reqBulkStockActivity *models.BulkAdjustment) ([]*schema.StockActivity, error) {
	err := reqBulkStockActivity.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	for _, stockActivityInstance := range reqBulkStockActivity.StockActivities {
		newStockActivity := getStockActivity(stockActivityInstance, requestParams)
		stock, err := sas.stockService.FindByID(requestParams, stockActivityInstance.StockID)
		if err != nil {
			return nil, err
		}

		_, err = sas.purposeService.FindByID(requestParams, stockActivityInstance.PurposeID)
		if err != nil {
			return nil, err
		}

		if stock.Status == consts.StockInactiveStatus {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.UnKnownErr,
				TranslationKey: "statusInactive",
				TranslationParams: map[string]interface{}{
					"field": "stock",
				},
				HttpCode: http.StatusBadRequest,
			}
		}

		stockActivity, err := sas.stockActivityRepository.CreateStockActivity(newStockActivity)
		if err != nil {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "createUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "Stock activity",
				},
				HttpCode: http.StatusInternalServerError,
			}
		}

		if stockActivityInstance.PurposeID != 0 {

			updatePurposeReq := &models.UpdatePurposeRequestBody{
				IsUsed: true,
			}
			_, err = sas.purposeService.UpdatePurpose(requestParams, stockActivity.PurposeID, updatePurposeReq)
			if err != nil {
				return nil, err
			}
		}

		newStock := getUpdatedStock(stockActivityInstance, stock)
		_, err = sas.stockService.UpdateStock(requestParams, stockActivityInstance.StockID, newStock)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
