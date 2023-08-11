package models

import (
	"pi-inventory/common/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (p Page) GetPageInformation(ctx *gin.Context) (*Page, error) {
	page := Page{Offset: 0, Limit: 10}
	var err error
	page.Offset, err = strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		if len(ctx.Query("offset")) != 0 {
			logger.LogError(err)
			return nil, err
		} else if len(ctx.Query("offset")) == 0 {
			page.Offset = 0
		}
	}
	page.Limit, err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		if len(ctx.Query("limit")) != 0 {
			logger.LogError(err)
			return nil, err
		} else if len(ctx.Query("limit")) == 0 {
			page.Limit = 10
		}
	}

	return &page, nil
}

type PageResponse struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
}
