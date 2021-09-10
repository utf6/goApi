package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/utf6/goApi/app"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/config"
	errors "github.com/utf6/goApi/pkg/error"
	"github.com/utf6/goApi/pkg/logger"
	"github.com/utf6/goApi/pkg/util"
	"net/http"
)

// @Tags 标签管理
// @Summary 获取文章标签
// @Param name query string false "标签名称"
// @Param state query int false "状态（0：禁用，1：正常）"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	//state := c.DefaultQuery("state", "1")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	name := c.Query("name")
	if name != "" {
		maps["name"] = name
	}

	state := 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	//获取数据
	data["list"] = models.GetTags(util.GetPage(c), config.Apps.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	app.Response(http.StatusOK, errors.SUCCESS, data, c)
}

// @Tags 标签管理
// @Summary 新增文章标签
// @Param name formData string true "标签名称"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.PostForm("name")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空！")
	valid.MaxSize(name, 100, "name").Message("名称最多100字符")

	//数据验证错误
	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断标签是否存在
	if !models.ExistTagByName(name) {
		app.Response(http.StatusBadRequest, errors.ERROR_EXIST_TAG, nil, c)
		return
	}

	//判断是否添加成功
	if models.AddTag(name) {
		//返回结果
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}

// @Tags 标签管理
// @Summary 修改文章标签
// @Param id path int true "ID"
// @Param token path string true "token"
// @Param name formData string true "标签名称"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")

	valid := validation.Validation{}
	valid.Required(id, "id").Message("id不能为空")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最多100字符")

	//数据验证错误
	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断标签是否存在
	if !models.ExistTagById(id) {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	//组合数据
	data := make(map[string]interface{})
	data["name"] = name
	if models.EditTag(id, data) {
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}

// @Tags 标签管理
// @Summary 删除文章标签
// @Param id path int true "ID"
// @Param token path string true "token"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [Delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id 必须大于0")

	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断标签是否存在
	if !models.ExistTagById(id) {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	if models.DeleteTag(id) {
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}