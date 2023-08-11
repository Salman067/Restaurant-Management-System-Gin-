package service

import (
	"net/http"
	"pi-inventory/common/logger"
	"pi-inventory/errors"
	attachmentService "pi-inventory/modules/attachment/service"
	compositeConst "pi-inventory/modules/composite/consts"
	"pi-inventory/modules/composite/models"
	"pi-inventory/modules/composite/repository"
	"pi-inventory/modules/composite/schema"
	"time"

	commonModel "pi-inventory/common/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CompositeServiceInterface interface {
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.CompositeQueryParams) ([]*models.SingleCompositeResponseBody, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.SingleCompositeResponseBody, error)
	CreateComposite(requestParams commonModel.AccountRequstParams, reqComposite *models.AddCompositeRequestBody) (*schema.Composite, error)
	UpdateComposite(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateCompositeRequestBody) (*schema.Composite, error)
	DeleteComposite(requestParams commonModel.AccountRequstParams, ID uint64) error
	CompositeItemSummary(requestParams commonModel.AccountRequstParams) (int64, error)
}

type compositeService struct {
	compositeRepository repository.CompositeRepositoryInterface
	lineItemService     LineItemServiceInterface
	attachmentService   attachmentService.AttachmentServiceInterface
}

func NewCompositeService(compositeRepo repository.CompositeRepositoryInterface, lineItemService LineItemServiceInterface, attachmetService attachmentService.AttachmentServiceInterface) *compositeService {
	return &compositeService{
		compositeRepository: compositeRepo,
		lineItemService:     lineItemService,
		attachmentService:   attachmetService,
	}
}
func (cs *compositeService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *models.CompositeQueryParams) ([]*models.SingleCompositeResponseBody, *commonModel.PageResponse, error) {
	composites, recordCount, err := cs.compositeRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		logger.LogError(err)
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	responseComposites, err := setCompositeResponse(composites)
	if err != nil {
		logger.LogError(err)
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for i, singleComposite := range composites {
		lineItems, err := cs.lineItemService.FetchLineItems(requestParams, singleComposite.LineItemKey)
		if err != nil {
			logger.LogError(err)
			return nil, nil, err
		}
		setlineItems(responseComposites[i], lineItems)
	}
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}
	return responseComposites, &pageInfo, nil
}

func (cs *compositeService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*models.SingleCompositeResponseBody, error) {
	singleComposite, err := cs.compositeRepository.FindBy(requestParams, compositeConst.CompositeFieldID, ID)
	if err != nil {
		logger.LogError(err)
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "composite item",
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
	responseComposites, err := setCompositeResponse(singleComposite)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	lineItems, err := cs.lineItemService.FetchLineItems(requestParams, singleComposite[0].LineItemKey)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	setlineItems(responseComposites[0], lineItems)
	attachments, err := cs.attachmentService.FetchAttachments(singleComposite[0].AttachmentKey)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	responseComposites[0].Attachments = append(responseComposites[0].Attachments, attachments...)
	return responseComposites[0], nil
}

func (cs *compositeService) CreateComposite(requestParams commonModel.AccountRequstParams, reqComposite *models.AddCompositeRequestBody) (*schema.Composite, error) {
	err := reqComposite.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	_, err = cs.compositeRepository.FindBy(requestParams, compositeConst.CompositeFieldTitle, reqComposite.Title)

	if err == nil {
		logger.LogError("A CompositeItem with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "composite item",
			},
			HttpCode: http.StatusBadRequest,
		}
	}

	lineItemKey, err := cs.lineItemService.CreateLineItem(requestParams, reqComposite.LineItems)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	newComposite := &schema.Composite{}
	setRequestCompositeFieldValues(newComposite, reqComposite, requestParams, lineItemKey)
	composite, err := cs.compositeRepository.Create(newComposite)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Composite item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return composite, nil
}

func (cs *compositeService) UpdateComposite(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *models.UpdateCompositeRequestBody) (*schema.Composite, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	composites, err := cs.compositeRepository.FindBy(requestParams, compositeConst.CompositeFieldID, ID)
	if err != nil {
		logger.LogError(err)
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "composite item",
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
	err = cs.lineItemService.DeleteLineItem(requestParams, composites[0].LineItemKey)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	newLineItemKey, err := cs.lineItemService.CreateLineItem(requestParams, reqBody.LineItems)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	checkedComposites := updateCompositeFields(reqBody, composites)
	checkedComposites[0].LineItemKey = newLineItemKey
	updatedComposite, err := cs.compositeRepository.UpdateComposite(checkedComposites)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Composite item",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	return updatedComposite, nil
}

func (cs *compositeService) DeleteComposite(requestParams commonModel.AccountRequstParams, ID uint64) error {
	panic("unimplemented")
}

func (cs *compositeService) CompositeItemSummary(requestParams commonModel.AccountRequstParams) (int64, error) {
	compositeItemCount, err := cs.compositeRepository.CompositeItemSummary(requestParams)
	return compositeItemCount, err
}

func setRequestCompositeFieldValues(newInstance *schema.Composite, resComposite *models.AddCompositeRequestBody, requestParams commonModel.AccountRequstParams, lineItemKey string) {
	newInstance.OwnerID = requestParams.CreatedBy
	newInstance.AccountID = requestParams.AccountID
	newInstance.Title = resComposite.Title
	newInstance.Tag = resComposite.Tag
	newInstance.Description = resComposite.Description
	newInstance.SellingPrice = resComposite.SellingPrice.Amount
	newInstance.SellingPriceCurrency = resComposite.SellingPrice.Currency
	newInstance.PurchasePrice = resComposite.PurchasePrice.Amount
	newInstance.PurchasePriceCurrency = resComposite.PurchasePrice.Currency
	newInstance.AttachmentKey = resComposite.AttachmentKey
	newInstance.CreatedAt = time.Now()
	newInstance.CreatedBy = requestParams.CreatedBy
	newInstance.LineItemKey = lineItemKey
}
func setResponseCompositeFieldValues(newInstance *schema.Composite, resComposite *models.SingleCompositeResponseBody) {
	resComposite.ID = newInstance.ID
	resComposite.OwnerID = newInstance.OwnerID
	resComposite.AccountID = newInstance.AccountID
	resComposite.Title = newInstance.Title
	resComposite.Description = newInstance.Description
	resComposite.Tag = newInstance.Tag
	resComposite.SellingPrice = commonModel.Money{
		Amount:   newInstance.SellingPrice,
		Currency: newInstance.SellingPriceCurrency,
	}.ConvertFloatMoney()
	resComposite.PurchasePrice = commonModel.Money{
		Amount:   newInstance.PurchasePrice,
		Currency: newInstance.PurchasePriceCurrency,
	}.ConvertFloatMoney()
	resComposite.AttachmentKey = newInstance.AttachmentKey
}

func setCompositeResponse(list []*schema.Composite) ([]*models.SingleCompositeResponseBody, error) {

	responseComposites := make([]*models.SingleCompositeResponseBody, 0)
	for _, singleInstance := range list {

		composite := models.SingleCompositeResponseBody{}
		setResponseCompositeFieldValues(singleInstance, &composite)
		responseComposites = append(responseComposites, &composite)
	}

	return responseComposites, nil
}
func setlineItems(responseComposite *models.SingleCompositeResponseBody, lineItems []*models.LineItemResponseBody) {
	responseComposite.LineItems = append(responseComposite.LineItems, lineItems...)
}

func updateCompositeFields(reqBody *models.UpdateCompositeRequestBody, composites []*schema.Composite) []*schema.Composite {
	if reqBody.Title != "" {
		composites[0].Title = reqBody.Title
	}
	if reqBody.Description != "" {
		composites[0].Description = reqBody.Description
	}
	if reqBody.SellingPrice.Amount.GreaterThan(decimal.NewFromFloat32(0)) {
		composites[0].SellingPrice = reqBody.SellingPrice.Amount
	}
	if reqBody.PurchasePrice.Amount.GreaterThan(decimal.NewFromFloat32(0)) {
		composites[0].PurchasePrice = reqBody.PurchasePrice.Amount
	}
	if reqBody.SellingPrice.Currency != "" {
		composites[0].SellingPriceCurrency = reqBody.SellingPrice.Currency
	}
	if reqBody.PurchasePrice.Currency != "" {
		composites[0].PurchasePriceCurrency = reqBody.PurchasePrice.Currency
	}
	return composites
}
