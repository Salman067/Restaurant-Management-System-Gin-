package service

import (
	"fmt"
	"net/http"
	"pi-inventory/common/cache"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	stockConst "pi-inventory/modules/stock/consts"
	stockModel "pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
	"pi-inventory/modules/stock/schema"
	"time"

	"gorm.io/gorm"
)

type CategoryServiceInterface interface {
	FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.CategoryQueryParams) ([]*stockModel.CategoryResponseBody, *commonModel.PageResponse, error)
	FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*stockModel.CategoryResponseBody, error)
	CreateCategory(requestParams commonModel.AccountRequstParams, reqCategory *stockModel.AddCategoryRequestBody) (*schema.Category, error)
	UpdateCategory(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *stockModel.UpdateCategoryRequestBody) (*schema.Category, error)
	DeleteCategory(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error)
}

type categoryService struct {
	categoryRepository repository.CategoryRepositoryInterface
	stockRepository    repository.StockRepositoryInterface
	cache              cache.CacheInterface
}

func NewCategoryService(categoryRepo repository.CategoryRepositoryInterface, stockRepository repository.StockRepositoryInterface, cache cache.CacheInterface) *categoryService {
	return &categoryService{
		categoryRepository: categoryRepo,
		stockRepository:    stockRepository,
		cache:              cache,
	}
}

func (cs *categoryService) FindAll(requestParams commonModel.AccountRequstParams, queryParamStruct *stockModel.CategoryQueryParams) ([]*stockModel.CategoryResponseBody, *commonModel.PageResponse, error) {
	var jsonResponse []*stockModel.CategoryResponseBody
	categories, recordCount, err := cs.categoryRepository.FindAll(requestParams, queryParamStruct)
	if err != nil {
		return nil, nil, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	for _, category := range *categories {
		resCategory := stockModel.CategoryResponseBody{
			ID:          category.ID,
			OwnerID:     category.OwnerID,
			AccountID:   category.AccountID,
			Title:       category.Title,
			Status:      category.Status,
			Description: category.Description,
			IsUsed:      category.IsUsed,
		}
		jsonResponse = append(jsonResponse, &resCategory)
	}
	pageInfo := commonModel.PageResponse{
		Offset: requestParams.Page.Offset,
		Limit:  requestParams.Page.Limit,
		Count:  int(recordCount),
	}
	return jsonResponse, &pageInfo, nil
}

func (cs *categoryService) FindByID(requestParams commonModel.AccountRequstParams, ID uint64) (*stockModel.CategoryResponseBody, error) {

	category, err := cs.categoryRepository.FindBy(requestParams, stockConst.CategoryFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "category",
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

	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Category, requestParams.CreatedBy, ID)
	err = cs.cache.Set(cacheKey, category.Title)
	if err != nil {
		logger.LogError(err)
	}

	resCategory := &stockModel.CategoryResponseBody{
		ID:          category.ID,
		OwnerID:     category.OwnerID,
		AccountID:   category.AccountID,
		Title:       category.Title,
		Status:      category.Status,
		Description: category.Description,
		IsUsed:      category.IsUsed,
	}
	return resCategory, nil
}

func (cs *categoryService) CreateCategory(requestParams commonModel.AccountRequstParams, reqCategory *stockModel.AddCategoryRequestBody) (*schema.Category, error) {
	err := reqCategory.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	_, err = cs.categoryRepository.FindBy(requestParams, stockConst.CategoryFieldTitle, reqCategory.Title)

	if err == nil {
		logger.LogError("A category with the title exists")
		return nil, &errors.ApplicationError{
			ErrorType:      errors.AlreadyExistsErr,
			TranslationKey: "alreadyExists",
			TranslationParams: map[string]interface{}{
				"field": "category",
			},
			HttpCode: http.StatusBadRequest}
	}
	newCategory := &schema.Category{
		Title:       reqCategory.Title,
		OwnerID:     requestParams.CreatedBy,
		Description: reqCategory.Description,
		IsUsed:      reqCategory.IsUsed,
		CreatedAt:   time.Now(),
		CreatedBy:   requestParams.CreatedBy,
		AccountID:   requestParams.AccountID,
	}
	category, err := cs.categoryRepository.Create(newCategory)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "createUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Category",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}
	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Category, requestParams.CreatedBy, category.ID)
	err = cs.cache.Set(cacheKey, category.Title)
	if err != nil {
		logger.LogError(err)
	}
	return category, nil
}

func (cs *categoryService) UpdateCategory(requestParams commonModel.AccountRequstParams, ID uint64, reqBody *stockModel.UpdateCategoryRequestBody) (*schema.Category, error) {
	err := reqBody.Validate()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	category, err := cs.categoryRepository.FindBy(requestParams, stockConst.CategoryFieldID, ID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "category",
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

	if reqBody.Title != "" {
		category.Title = reqBody.Title
	}
	if reqBody.Status != "" {
		category.Status = reqBody.Status
	}
	if reqBody.Description != "" {
		category.Description = reqBody.Description
	}
	category.AccountID = requestParams.AccountID
	category.IsUsed = reqBody.IsUsed
	category.UpdatedAt = time.Now()
	category.OwnerID = requestParams.CreatedBy

	category, err = cs.categoryRepository.Update(category)
	if err != nil {
		return nil, &errors.ApplicationError{
			ErrorType:      errors.Unsuccessfull,
			TranslationKey: "updateUnsuccessfull",
			TranslationParams: map[string]interface{}{
				"field": "Category",
			},
			HttpCode: http.StatusInternalServerError,
		}
	}

	cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Category, requestParams.CreatedBy, category.ID)
	err = cs.cache.Set(cacheKey, category.Title)
	if err != nil {
		logger.LogError(err)
	}
	return category, nil
}

func (cs *categoryService) DeleteCategory(requestParams commonModel.AccountRequstParams, ID uint64) (uint64, error) {

	_, err := cs.categoryRepository.FindBy(requestParams, stockConst.CategoryFieldID, ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.RecordNotFound,
				TranslationKey: "idNotFound",
				TranslationParams: map[string]interface{}{
					"field": "category",
				},
				HttpCode: http.StatusBadRequest,
			}
		}
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownErr",
			HttpCode:       http.StatusInternalServerError,
		}
	}
	//Check Before Delete
	stocks, err := cs.stockRepository.FindBy(requestParams, stockConst.StockFieldCategoryID, ID)
	if err != nil {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "unknownError",
			HttpCode:       http.StatusInternalServerError,
		}
	} else if len(stocks) > 0 {
		return 0, &errors.ApplicationError{
			ErrorType:      errors.UnKnownErr,
			TranslationKey: "deleteFailedDueToDependency",
			TranslationParams: map[string]interface{}{
				"field": "category",
			},
			HttpCode: http.StatusBadRequest,
		}
	} else {
		categoryID, err := cs.categoryRepository.Delete(requestParams, ID)
		if err != nil {
			return 0, &errors.ApplicationError{
				ErrorType:      errors.Unsuccessfull,
				TranslationKey: "deleteUnsuccessfull",
				TranslationParams: map[string]interface{}{
					"field": "Category",
				},
				HttpCode: http.StatusInternalServerError}
		}
		cacheKey := fmt.Sprintf("%s:%d:%d", stockConst.Category, requestParams.CreatedBy, ID)
		err = cs.cache.Remove(cacheKey)
		if err != nil {
			logger.LogError(err)
		}
		return categoryID, nil
	}

}
