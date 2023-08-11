package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pi-inventory/common/consts"
	"pi-inventory/common/logger"
	commonModel "pi-inventory/common/models"
	"pi-inventory/errors"
	"pi-inventory/modules/profile/models"
	"pi-inventory/modules/profile/service"
)

type ProfileControllerInterface interface {
	Pong(context *gin.Context)
	ProfileView(context *gin.Context)
}

type profileController struct {
	errors         errors.GinError
	profileService service.ProfileServiceInterface
}

func NewProfileController(profileSrv service.ProfileServiceInterface) *profileController {
	return &profileController{profileService: profileSrv}
}

func (pc *profileController) Pong(context *gin.Context) {
	context.JSON(http.StatusOK, "Pong from Warehouse")
}

func (pc *profileController) ProfileView(context *gin.Context) {
	//ctx := context.Request.Context()
	//accountSlug := context.Params.ByName("account_slug")
	accountId := context.GetUint64("account_id")
	var accountInfo commonModel.RedisAccountInfo
	value, ok := context.Get("account_info")
	if ok {
		accountInfo = value.(commonModel.RedisAccountInfo)
	}
	if consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessSubBook || consts.AccountTypes[accountInfo.Type] == consts.AccountTypeBusinessBranchBook {
		accountId = uint64(accountInfo.MainAccountId)
		//accountSlug = accountInfo.MainAccountSlug
	}

	resProfile := models.Profile{}
	ownerID := context.GetInt64("user_id")
	for _, accountUser := range accountInfo.AccountUserPermissions {
		logger.LogError(accountUser.UserId, " ", ownerID)
		if accountUser.UserId == uint(ownerID) {
			profile := models.Profile{
				ID:             accountUser.User.ID,
				AccountID:      accountId,
				Name:           accountUser.User.Name,
				Email:          accountUser.User.Email,
				Mobile:         accountUser.User.Mobile,
				ProfilePicture: accountUser.User.ProfilePicture,
			}
			resProfile = profile
		}
	}
	context.JSON(http.StatusOK, gin.H{"profile": resProfile})
}
