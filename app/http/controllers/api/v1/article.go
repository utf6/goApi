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

// @Tags 文章管理
// @Summary 获取单个文章
// @Param id path int false "文章id"
// @Param token path string true "access_token"
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("非法id")

	//数据验证错误
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logger.Error(err.Key, err.Message)
		}
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断文章是否存在
	if models.ExistArticleByID(id) {
		app.Response(http.StatusOK, errors.SUCCESS, models.GetArticle(id), c)
		return
	}

	app.Response(http.StatusOK, errors.ERROR_NOT_EXIST_ARTICLE, models.GetArticle(id), c)
}

// @Tags 文章管理
// @Summary 获取多个文章
// @Param tag_id query int false "标签id"
// @Param state query int false "状态（0：删除，1：正常）"
// @Param token path string true "access_token"
// @Failure 400 {object} app.Response
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	if state := com.StrTo(c.Query("state")).MustInt() ; state > -1 {
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
		maps["state"] = state
	}

	if tagId := com.StrTo(c.Query("tag_id")).MustInt(); tagId > 0 {
		valid.Min(tagId, 1, "tag_id").Message("标签id错误")
		maps["tag_id"] = tagId
	}

	//验证错误
	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = models.GetArticles(util.GetPage(c), config.Apps.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)
	app.Response(http.StatusOK, errors.SUCCESS, data, c)
}

// @Tags 文章管理
// @Summary 新增文章
// @Param tag_id formData int true "标签id"
// @Param title formData string true "文章标题"
// @Param desc formData string true "文章描述"
// @Param content formData string true "文章内容"
// @Param token path string true "access_token"
// @Failure 500 {object} app.Response
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	thumb := c.PostForm("thumb")
	desc := c.PostForm("desc")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(tagId, "tag_id").Message("标签不能为空")
	valid.Min(tagId, 1, "tag_id").Message("标签错误")
	valid.Required(title, "title").Message("标题不能为空")
	valid.MaxSize(title, 100, "title").Message("标题最多100个字符")
	valid.MaxSize(desc, 200, "desc").Message("描述最多200字")
	valid.Required(content, "content").Message("内容不能为空")

	//判断图片是否存在
	//if thumb != "" && !com.IsExist(thumb) {
	//	app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil, c)
	//	return
	//}

	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	if !models.ExistTagById(tagId) {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	//插入数据
	data := make(map[string]interface{})
	data["thumb"] = thumb
	data["title"] = title
	data["desc"] = desc
	data["tag_id"] = tagId
	data["content"] = content
	if models.AddArticle(data) {
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}

// @Tags 文章管理
// @Summary 修改文章
// @Param id path int true "文章id"
// @Param tag_id formData int true "标签id"
// @Param title formData string true "文章标题"
// @Param desc formData string true "文章描述"
// @Param content formData string true "文章内容"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	//组合数据
	id := com.StrTo(c.Param("id")).MustInt()
	title := c.PostForm("title")
	thumb := c.PostForm("thumb")
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

	//判断图片是否存在
	//if thumb != "" && !com.IsExist(thumb) {
	//	app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil, c)
	//	return
	//}

	//验证失败
	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断标签是否存在
	if !models.ExistTagById(tagId) {
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//派单文章是否存在
	if !models.ExistArticleByID(id) {
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//组合数据
	data := make(map[string]interface{})
	data["thumb"] = thumb
	data["title"] = title
	data["desc"] = desc
	data["tag_id"] = tagId
	data["content"] = content
	//修改数据
	if models.EditArticle(id, data) {
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	//返回错误
	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}

// @Tags 文章管理
// @Summary 删除文章
// @Param id path int true "文章id"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [Delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	//判断文章是否存在
	if !models.ExistArticleByID(id) {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_ARTICLE, nil, c)
		return
	}

	//删除数据
	if models.DeleteArticle(id) {
		app.Response(http.StatusOK, errors.SUCCESS, nil, c)
		return
	}

	//返回错误
	app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
}