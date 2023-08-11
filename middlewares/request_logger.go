package middlewares

import (
	"bytes"
	"io"
	"pi-inventory/common/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		inproccessURLs := viper.GetString("RQUEST_BODY_AVOID_URLS")
		proccessedURLs := strings.Split(inproccessURLs, ",")
		for _, url := range proccessedURLs {
			if url == c.Request.URL.Path {
				return
			}
		}

		if viper.GetString("GIN_MODE") == "release" {
			c.Next()
			return
		}
		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		body := readBody(rdr1)
		path := c.Request.URL.Path
		queryParams := c.Request.URL.Query()
		logger.LogDebug("path ", path)
		logger.LogDebug("query params ", queryParams)
		logger.LogDebug("rquest body ", body)

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
