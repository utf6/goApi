package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/config"
	errors "github.com/utf6/goApi/pkg/error"
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
	name := c.Query("name")
	//state := c.DefaultQuery("state", "1")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	state := 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := errors.SUCCESS
	data["list"] = models.GetTags(util.GetPage(c), config.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": data,
	})
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

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		//判断标签是否存在
		if !models.ExistTagByName(name) {
			//判断是否添加成功
			if models.AddTag(name) {
				code = errors.SUCCESS
			} else {
				code = errors.ERROR
			}
		} else {
			code = errors.ERROR_EXIST_TAG
		}
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]string),
	})

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

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(id) {
			data := make(map[string]interface{})
			data["name"] = name

			if models.EditTag(id, data) {
				code = errors.SUCCESS
			} else {
				code = errors.ERROR
			}
		} else {
			code = errors.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]string),
	})
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

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(id) {
			if models.DeleteTag(id) {
				code = errors.SUCCESS
			} else {
				code = errors.ERROR
			}
		} else {
			code = errors.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]string),
	})
}
