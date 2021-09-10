package facade

import "C"
import (
	"github.com/gin-gonic/gin"
	errors "github.com/utf6/goApi/pkg/error"
)

type Response struct {

	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func responses(httpCode, code int, data interface{}, C *gin.Context)  {
	C.JSON(httpCode, gin.H{
		"code" : code,
		"msg" : errors.GetMsg(code),
		"data": data,
	})
}
