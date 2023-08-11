package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shopspring/decimal"
	"github.com/ttacon/libphonenumber"
)

func Reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func CopyStructToStruct(input interface{}, output interface{}) error {
	if byteData, err := json.Marshal(input); err == nil {
		if err := json.Unmarshal(byteData, &output); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}

func ParsingContactNumber(number string, numberCountry string) (string, error) {
	num, err := libphonenumber.Parse(number, strings.ToUpper(numberCountry))
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("unKnownError")
	}

	var valid = libphonenumber.IsValidNumberForRegion(num, strings.ToUpper(numberCountry))
	if !valid {
		return "", errors.New("invalid")
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	return formattedNum, nil
}
func CheckValidContactNumber(number string, numberCountry string) error {
	_, err := ParsingContactNumber(number, numberCountry)

	return err
}

func ScalePrice(unitPrice decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(unitPrice.InexactFloat64()).RoundUp(2)
}

func StringToDate(date string) (time.Time, error) {
	formattedDate, err := time.Parse(time.RFC3339, date)
	return formattedDate, err
}

func Param(c *gin.Context) (uint64, error) {
	tempID := c.Param("id")
	ID, err := strconv.ParseUint(tempID, 0, 0)
	if err != nil && tempID != "" {
		return 0, err
	}
	return ID, nil
}

// func GetOwnerIDFromCtx(context *gin.Context) int64 {
// 	ownerID := context.GetInt64("user_id")

// 	return ownerID
// }
