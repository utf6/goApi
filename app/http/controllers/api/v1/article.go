package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/utf6/goApi/app"
	"github.com/utf6/goApi/app/repository"
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
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//数据仓库
	articleRepository := repository.Article{ID:id}
	article := articleRepository.ExistByID()

	//判断文章是否存在
	if !article {
		app.Response( http.StatusOK, errors.ERROR_NOT_EXIST_ARTICLE, nil, c)
		return
	}

	//获取数据
	result, err := articleRepository.Get()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_GET_ARTICLE_FAIL, nil, c)
		return
	}

	app.Response(http.StatusOK, errors.SUCCESS, result, c)
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
	valid := validation.Validation{}

	 state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()
	 if state >= 0 {
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}

	 tagId := com.StrTo(c.Query("tag_id")).MustInt()
	 if tagId > 0 {
		valid.Min(tagId, 1, "tag_id").Message("标签id错误")
	}

	//验证错误
	if valid.HasErrors() {
		logger.Errors(valid.Errors)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//获取数据仓库
	articleRepository := repository.Article{
		TagID:     tagId,
		State:     state,
		PageNum:   util.GetPage(c),
		PageSize:  config.Apps.PageSize,
	}

	//获取文章总数
	total, err := articleRepository.Count()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_COUNT_ARTICLE_FAIL, nil, c)
		return
	}

	//获取所有文章
	articles, err := articleRepository.GetAll()
	if err != nil {
		app.Response(http.StatusNonAuthoritativeInfo, errors.ERROR_GET_ARTICLE_FAIL, nil, c)
		return
	}

	//返回数据
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
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

	//判断标签是否存在
	tagRepository := repository.Tag{ID:tagId}
	exist := tagRepository.ExistById()
	if !exist {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	//生成文章数据仓库
	articleRepository := repository.Article{
		TagID:     tagId,
		Title:     title,
		Desc:      desc,
		Content:   content,
		Thumb:     thumb,
		State:     0,
	}
	if err := articleRepository.Add(); err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_ADD_ARTICLE_FAIL, nil, c)
	}
	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
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
	tagRepository := repository.Tag{ID:tagId}
	tag := tagRepository.ExistById()
	if !tag {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	//生成文章数据仓库
	articleRepository := repository.Article{
		ID:        id,
		TagID:     tagId,
		Title:     title,
		Desc:      desc,
		Content:   content,
		Thumb:     thumb,
		State:     1,
	}

	//判断文章是否存在
	article := articleRepository.ExistByID()
	if !article {
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//修改数据
	if err := articleRepository.Edit(); err != nil {
		//返回错误
		app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
		return
	}
	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
}

// @Tags 文章管理
// @Summary 删除文章
// @Param id path int true "文章id"
// @Param token path string true "access_token"
// @Success 200 {object} gin.H "{"code":200, "data":{}, "msg":"ok"}"
// @Router /api/v1/articles/{id} [Delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	//生成文章数据仓库
	articleRepository := repository.Article{ID:id}

	//判断文章是否存在
	article := articleRepository.ExistByID()
	if !article {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_ARTICLE, nil, c)
		return
	}

	//删除数据
	if err := articleRepository.Delete(); err != nil {
		//返回错误
		app.Response(http.StatusInternalServerError, errors.ERROR, nil, c)
		return
	}

	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
}