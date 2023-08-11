package utils

import (
	"fmt"
	"strconv"
)

func Num64(n interface{}) int64 {
	s := fmt.Sprintf("%d", n)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}
