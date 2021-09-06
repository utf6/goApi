package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/config"
	errors "github.com/utf6/goApi/pkg/error"
	"github.com/utf6/goApi/pkg/logger"
	"github.com/utf6/goApi/pkg/util"
	"net/http"
)

// @Tags 文章管理
// @Summary 获取单个文章
// @Param id path int false "文章id"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("非法id")

	code := errors.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = errors.SUCCESS
		} else {
			code = errors.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logger.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": data,
	})
}

// @Tags 文章管理
// @Summary 获取多个文章
// @Param tag_id query int false "标签id"
// #Param state query int false "状态（0：删除，1：正常）"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	if state := c.Query("state"); state != "" {
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
		maps["state"] = state
	}

	if tagId := c.Query("tag_id"); tagId != "" {
		valid.Min(tagId, 1, "tag_id").Message("标签id错误")
		maps["tag_id"] = tagId
	}

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		code = errors.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), config.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logger.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": data,
	})
}

// @Tags 文章管理
// @Summary 新增文章
// @Param tag_id formData int true "标签id"
// #Param title formData string true "文章标题"
// #Param desc formData string true "文章描述"
// #Param content formData string true "文章内容"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(tagId, "tag_id").Message("标签不能为空")
	valid.Min(tagId, 1, "tag_id").Message("标签错误")
	valid.Required(title, "title").Message("标题不能为空")
	valid.MaxSize(title, 100, "title").Message("标题最多100个字符")
	valid.MaxSize(desc, 200, "desc").Message("描述最多200字")
	valid.Required(content, "content").Message("内容不能为空")

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(tagId) {
			data := make(map[string]interface{})
			data["title"] = title
			data["desc"] = desc
			data["tag_id"] = tagId
			data["content"] = content
			if models.AddArticle(data) {
				code = errors.SUCCESS
			} else {
				code = errors.ERROR
			}
		} else {
			code = errors.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logger.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Tags 文章管理
// @Summary 修改文章
// @Param id path int true "文章id"
// @Param tag_id formData int true "标签id"
// #Param title formData string true "文章标题"
// #Param desc formData string true "文章描述"
// #Param content formData string true "文章内容"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	//组合数据
	id := com.StrTo(c.Param("id")).MustInt()
	title := c.PostForm("title")
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	desc := c.PostForm("desc")
	content := c.PostForm("content")

	//验证数据
	valid := validation.Validation{}
	valid.Required(id, "id").Message("参数错误")
	valid.Required(title, "title").Message("标题不能为空")
	valid.MaxSize(title, 100, "title").Message("标题最多100个字符")
	valid.MaxSize(desc, 200, "desc").Message("描述不能为空")
	valid.Required(content, "content").Message("内容不能为空")

	code := errors.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(tagId) {
			data := make(map[string]interface{})
			data["title"] = title
			data["desc"] = desc
			data["tag_id"] = tagId
			data["content"] = content

			if models.EditArticle(id, data) {
				code = errors.SUCCESS
			} else {
				code = errors.ERROR
			}
		} else {
			code = errors.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logger.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// @Tags 文章管理
// @Summary 删除文章
// @Param id path int true "文章id"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [Delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := errors.INVALID_PARAMS
	if models.ExistArticleByID(id) {
		if models.DeleteArticle(id) {
			code = errors.SUCCESS
		} else {
			code = errors.ERROR
		}
	} else {
		code = errors.ERROR_NOT_EXIST_ARTICLE
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errors.GetMsg(code),
		"data": make(map[string]string),
	})
}
