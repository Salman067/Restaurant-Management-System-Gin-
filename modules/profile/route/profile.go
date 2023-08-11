package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	"pi-inventory/dic"
	"pi-inventory/middlewares"
	profileConst "pi-inventory/modules/profile/consts"
	profileController "pi-inventory/modules/profile/controller"
)

func SetupProfileRoute(_ *di.Builder, api *gin.RouterGroup) {
	profileController := dic.Container.Get(profileConst.ProfileController).(profileController.ProfileControllerInterface)
	profile := api.Group("/profile/:account_slug")
	profile.Use(middlewares.Auth(), middlewares.AccountDependValidate())
	{
		profile.GET("/ping", profileController.Pong)
		profile.GET("/view", profileController.ProfileView)
	}
}
