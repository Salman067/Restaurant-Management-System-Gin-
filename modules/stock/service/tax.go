package service

import (
	"context"
	"encoding/json"
	"pi-inventory/common/logger"
	"pi-inventory/modules/stock/cache"
	"pi-inventory/modules/stock/models"
	"pi-inventory/modules/stock/repository"
)

type TaxServiceInterface interface {
	FindAll(requestParams models.RequstParams) ([]*models.TaxResponseBody, error)
	FindByID(requestParams models.RequstParams, ID string) (*models.TaxResponseBody, error)
}
type taxService struct {
	taxRepository      repository.TaxRepositoryInterface
	taxCacheRepository cache.TaxCacheRepositoryInterface
}

func NewTaxService(TaxRepo repository.TaxRepositoryInterface, taxCacheRepository cache.TaxCacheRepositoryInterface) *taxService {
	return &taxService{taxRepository: TaxRepo, taxCacheRepository: taxCacheRepository}
}

func (ts *taxService) FindAll(requestParams models.RequstParams) ([]*models.TaxResponseBody, error) {

	key := "tax_" + requestParams.AccountSlug // Replace with your desired key pattern

	taxListMap, err := ts.taxCacheRepository.GetAll(context.Background(), key)
	if err != nil {
		logger.LogError(err)
	}
	taxRespList := make([]*models.TaxResponseBody, 0)

	for _, val := range taxListMap {
		var taxResp models.TaxResponseBody
		err := json.Unmarshal([]byte(val), &taxResp)
		if err != nil {
			logger.LogError("Error while unmarshal the json in tax response")
		}
		taxRespList = append(taxRespList, &taxResp)
	}

	return taxRespList, nil
}

func (ts *taxService) FindByID(requestParams models.RequstParams, ID string) (*models.TaxResponseBody, error) {

	key := "tax_" + requestParams.AccountSlug // Replace with your desired key pattern
	data, err := ts.taxCacheRepository.GetSingleData(context.Background(), key, ID)
	if err != nil {
		logger.LogError(err)
	}

	var taxResp models.TaxResponseBody
	err = json.Unmarshal([]byte(data), &taxResp)
	if err != nil {
		logger.LogError("Error while unmarshal the json in tax response")
	}

	return &taxResp, nil
}

// func demoTax(ID string) *models.TaxResponseBody {
// 	demoTaxList := db.DemoTaxList
// 	for _, val := range demoTaxList {
// 		if val.ID.String() == ID {
// 			return val
// 		}
// 	}
// 	return &models.TaxResponseBody{
// 		ID:          uuid.New(),
// 		TaxName:     "Demo Tax 1",
// 		AgencyID:    uuid.New(),
// 		Description: "Description",
// 		SalesRateHistory: []*models.RespRateHistory{
// 			{
// 				Rate:      10,
// 				StartDate: "",
// 			},
// 			{
// 				Rate:      20,
// 				StartDate: "",
// 			},
// 		},
// 		PurchaseRateHistory: []*models.RespRateHistory{
// 			{
// 				Rate:      5,
// 				StartDate: "",
// 			},
// 			{
// 				Rate:      6,
// 				StartDate: "",
// 			},
// 		},
// 		Status:    "active",
// 		CreatedBy: 1,
// 		AccountID: 1,
// 	}
// }

// func getDemoTaxList() []*models.TaxResponseBody {
// 	return db.DemoTaxList
// }
