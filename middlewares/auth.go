package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var accessToken string
		// cookie, err := ctx.Cookie("access_token")

		// authorizationHeader := ctx.Request.Header.Get("Authorization")
		// fields := strings.Fields(authorizationHeader)

		// if len(fields) != 0 && fields[0] == "Bearer" {
		// 	accessToken = fields[1]
		// } else if err == nil {
		// 	accessToken = cookie
		// }

		// if accessToken == "" {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
		// 	return
		// }

		// sub, err := utils.ValidateToken(accessToken, config.Config.AccessTokenPublicKey)
		// if err != nil {
		// 	logger.LogError(err)
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		// 	return
		// }

		// ctx.Set("user_id", int64(sub.(float64)))
		ctx.Set("user_id", int64(2))
		ctx.Next()
	}
}
