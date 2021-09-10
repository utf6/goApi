package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/app/models"
	errors "github.com/utf6/goApi/pkg/error"
	"github.com/utf6/goApi/pkg/files"
	"github.com/utf6/goApi/pkg/logger"
	"github.com/utf6/goApi/pkg/util"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Tags 系统管理
// @Summary 获取token
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/auth/getToken [post]
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
			} else {
				code = errors.SUCCESS
			}
		} else {
			code = errors.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logger.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": data,
	})
}

//@Tags 文件上传
//@Summary 获取token
//@Param username formData string true "用户名"
//@Param password formData string true "密码"
//@Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
//@Router /api/auth/getToken [post]
func Uploads(c *gin.Context)  {
	code := errors.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("files")
	if err != nil || image == nil {
		code = errors.INVALID_PARAMS
	} else {
		imageName := files.GetImageName(image.Filename)
		savePath := files.GetImagePath()

		src := "public/" + savePath + imageName
		if !files.CheckImageExt(imageName) || ! files.CheckImageSize(file) {
			code = errors.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := files.CheckImage("public/" + savePath)
			if err != nil {
				 logger.Warn(err)
				 code = errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logger.Warn(err)
				code = errors.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["file_path"] = files.GetImageFullUrl(imageName)
				data["file_name"] = imageName
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : errors.GetMsg(code),
		"data" : data,
	})
}