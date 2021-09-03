package middleware

import (
	"github.com/gin-gonic/gin"
	errors "github.com/utf6/goApi/pkg/error"
	"github.com/utf6/goApi/pkg/util"
	"net/http"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = errors.SUCCESS
		token := c.Query("token")

		if token == "" {
			code = errors.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)

			if err != nil {
				code = errors.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if claims.ExpiresAt < time.Now().Unix() {
				code = errors.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != errors.SUCCESS {
			 c.JSON(http.StatusUnauthorized, gin.H{
			 	"code" : code,
			 	"msg" : errors.GetMsg(code),
			 	"data" : data,
			 })
			 c.Abort()
			 return
		}
		c.Next()
	}
}
