package service

import (
	"net/http"
	"pi-inventory/common/logger"
	"pi-inventory/common/models"
	"pi-inventory/errors"
	compositeService "pi-inventory/modules/composite/service"
	groupItemService "pi-inventory/modules/groupItem/service"
)

type CommonServiceInterface interface {
	GroupItemAndCompositeItemSummary(requestParams models.AccountRequstParams) (*models.GroupItemAndCompositeItemSummaryResponse, error)
}

type commonService struct {
	compositeItemService compositeService.CompositeServiceInterface
	groupItemService     groupItemService.GroupItemServiceInterface
}

func NewCommonService(compositeItemService compositeService.CompositeServiceInterface,
	groupItemService groupItemService.GroupItemServiceInterface) *commonService {
	return &commonService{
		compositeItemService: compositeItemService,
		groupItemService:     groupItemService,
	}
}

func (cs *commonService) GroupItemAndCompositeItemSummary(requestParams models.AccountRequstParams) (*models.GroupItemAndCompositeItemSummaryResponse, error) {
	groupItemCount, err := cs.groupItemService.GroupItemSummary(requestParams)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{ErrorType: errors.UnKnownErr, TranslationKey: "unKnownFrr", HttpCode: http.StatusBadRequest}
	}

	compositeItemCount, err := cs.compositeItemService.CompositeItemSummary(requestParams)
	if err != nil {
		logger.LogError(err)
		return nil, &errors.ApplicationError{ErrorType: errors.UnKnownErr, TranslationKey: "unKnownFrr", HttpCode: http.StatusBadRequest}
	}

	response := &models.GroupItemAndCompositeItemSummaryResponse{
		GroupItemCount:     groupItemCount,
		CompositeItemCount: compositeItemCount,
	}
	return response, nil
}
