package response

// import (
// 	"fmt"
// 	"net/http"

// 	"pi-invoice/common/logger"
// 	"pi-invoice/rest_errors"
// )

// type Body map[string]interface{}

// func ValidationErrors(err error, entity string) (int, Body) {
// 	message := fmt.Sprintf("failed to validate the fields of the %v", entity)
// 	return validationResponse(message, err)
// }

// func GenerateErrorStatusCode(err error) (int) {
// 	message := err.Error()
// 	return readFromMap(message)
// }

// func readFromMap(message string) (int) {
// 	httpStatus, available := rest_errors.ResponseMap()[message]
// 	logger.LogError(message, " ",httpStatus)
// 	if available {
// 		return httpStatus
// 	}
// 	return http.StatusInternalServerError
// }

// func GenerateResponseBody(message string) Body {
// 	return Body{
// 		"message": message,
// 	}
// }

// func validationResponse(message string, err error) (int, Body) {
// 	return http.StatusBadRequest, Body{
// 		"message":          message,
// 		"validation_error": err,
// 	}
// }
