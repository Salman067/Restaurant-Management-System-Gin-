package controller

import (
	"encoding/json"
	"net/http"
	"pi-inventory/common/logger"
	"pi-inventory/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const StatusOk = "ok"
const StatusError = "error"

type BaseControllerInterface interface {
}

type BaseController struct {
	logger logger.LoggerInterface
}

func NewBaseController(logger logger.LoggerInterface) *BaseController {
	return &BaseController{logger: logger}
}

func (c BaseController) ReplySuccess(context *gin.Context, data interface{}) {
	c.Response(context, gin.H{"data": data, "status": StatusOk}, http.StatusOK)
}

func (c BaseController) SuccessResponse(context *gin.Context, data interface{}, data2 gin.H) {
	d := gin.H{"data": data, "status": StatusOk}
	if len(data2) > 0 {
		for k, v := range data2 {
			d[k] = v
		}
	}
	c.Response(context, d, http.StatusOK)
}

func (c BaseController) ReplyError(context *gin.Context, message string, code int) {
	c.Response(context, gin.H{"message": message, "status": StatusError}, code)
}

func (c BaseController) ReplyErrorWithType(context *gin.Context, message string, code int, errorType string) {
	c.Response(context, gin.H{"message": message, "status": StatusError, "type": errorType}, code)
}

func (c BaseController) Response(context *gin.Context, obj interface{}, code int) {
	switch context.GetHeader("Accept") {
	case "application/xml":
		context.XML(code, obj)
	default:
		context.JSON(code, obj)
	}
}

func (c BaseController) ReplyValidationError(context *gin.Context, err error) {
	errs := err.(validator.ValidationErrors)
	allErr := utils.TransValidationErrors(errs)
	context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "failed", "type": "validation_error", "data": allErr})
}

func GetAsString(item interface{}) string {
	out, err := json.Marshal(item)
	if err != nil {
		return ""
	}
	return string(out)
}
