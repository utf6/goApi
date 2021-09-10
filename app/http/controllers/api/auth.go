package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/app"
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

	//验证数据
	if !ok {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//验证用户
	isExist := models.CheckAuth(username, password)
	if !isExist {
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	var err error
	data := make(map[string]interface{})
	data, err = util.GenerateToken(username)
	if err != nil {
		app.Response(http.StatusBadRequest, errors.ERROR_AUTH_TOKEN, nil, c)
		return
	}

	app.Response(http.StatusOK, errors.SUCCESS, data, c)
}

//@Tags 文件上传
//@Summary 获取token
//@Param username formData string true "用户名"
//@Param password formData string true "密码"
//@Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
//@Router /api/auth/getToken [post]
func Uploads(c *gin.Context)  {
	file, image, err := c.Request.FormFile("files")
	if err != nil || image == nil {
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	imageName := files.GetImageName(image.Filename)
	savePath := files.GetImagePath()
	src := "public/" + savePath + imageName

	if !files.CheckImageExt(imageName) || ! files.CheckImageSize(file) {
		app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil, c)
		return
	}

	err = files.CheckImage("public/" + savePath)
	if err != nil {
		logger.Warn(err)
		app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil, c)
		return
	}

	err = c.SaveUploadedFile(image, src);
	if  err != nil {
		logger.Warn(err)
		app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil, c)
		return
	}

	data := make(map[string]string)
	data["file_path"] = files.GetImageFullUrl(imageName)
	data["file_name"] = imageName
	app.Response(http.StatusOK, errors.SUCCESS, data, c)
}