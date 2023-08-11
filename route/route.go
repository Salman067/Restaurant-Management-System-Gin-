package route

import (
	"net/http"
	commonConst "pi-inventory/common/consts"
	"pi-inventory/dic"
	"pi-inventory/middlewares"

	ginI18n "github.com/gin-contrib/i18n"

	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	commonModule "pi-inventory/common/route"
	attachmentModule "pi-inventory/modules/attachment/route"
	compositeModule "pi-inventory/modules/composite/route"
	groupItemModule "pi-inventory/modules/groupItem/route"
	profileModule "pi-inventory/modules/profile/route"
	stockModule "pi-inventory/modules/stock/route"
	warehouseModule "pi-inventory/modules/warehouse/route"
)

func Setup(_ *di.Builder) *gin.Engine {
	gin.SetMode(viper.GetString("GIN_MODE"))

	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(middlewares.SetTraceID())

	//Set localization
	r.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         "./localize",
		AcceptLanguage:   []language.Tag{language.Bengali, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    yaml.Unmarshal,
		FormatBundleFile: "yaml",
	})))
	commonConst.IsGinInitialized = true

	setupCors(r)

	// userValidation.InitCustomValidationRule()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(middlewares.LogAccessLog())
	r.Use(middlewares.RequestLogger())
	api := r.Group("/api")

	commonModule.SetupCommonRoute(dic.CommonBuilder, api)
	attachmentModule.SetupAttachmentRoute(dic.CommonBuilder, api)
	stockModule.SetupUnitRoute(dic.CommonBuilder, api)
	stockModule.SetupCategoryRoute(dic.CommonBuilder, api)
	stockModule.SetupStockRoute(dic.CommonBuilder, api)
	stockModule.SetupPurposeRoute(dic.CommonBuilder, api)
	stockModule.SetupStockActivityRoute(dic.CommonBuilder, api)
	stockModule.SetupTaxRoute(dic.CommonBuilder, api)
	stockModule.SetupSupplierRoute(dic.CommonBuilder, api)
	groupItemModule.SetupVariantRoute(dic.CommonBuilder, api)
	groupItemModule.SetupGroupItemRoute(dic.CommonBuilder, api)
	warehouseModule.SetupWarehouseRoute(dic.CommonBuilder, api)
	compositeModule.SetupCompositeRoute(dic.CommonBuilder, api)
	profileModule.SetupProfileRoute(dic.CommonBuilder, api)
	return r
}

func setupCors(r *gin.Engine) {
	allowConf := viper.GetString("CORS_ALLOW_ORIGINS")
	if allowConf == "" {
		r.Use(cors.Default())
		return
	}
	allowSites := strings.Split(allowConf, ",")
	config := cors.DefaultConfig()
	config.AllowOrigins = allowSites
}
