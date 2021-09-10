package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/app"
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
			app.Response(http.StatusUnauthorized, code, data, c)
			c.Abort()
			return
		}
		c.Next()
	}
}
