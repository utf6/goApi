package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/app/models"
	errors "github.com/utf6/goApi/pkg/error"
	"github.com/utf6/goApi/pkg/util"
	"log"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetToken(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid := validation.Validation{}
	ok, _ := valid.Valid(auth{
		Username: username,
		Password: password,
	})

	data := make(map[string]interface{})
	code := errors.INVALID_PARAMS

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			var err error
			data, err = util.GenerateToken(username)
			if err != nil {
				code = errors.ERROR_AUTH_TOKEN
			} else  {
				code = errors.SUCCESS
			}
		} else {
			code = errors.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : errors.GetMsg(code),
		"data" : data,
	})
}