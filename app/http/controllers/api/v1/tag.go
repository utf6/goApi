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
	"path"
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
	state := 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	//生成数据仓库
	tagRepository := repository.Tag{
		Name:      name,
		State:     state,
		PageNum:   util.GetPage(c),
		PageSize:  config.Apps.PageSize,
	}

	//获取所有
	list, err := tagRepository.GetAll()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_GET_TAG_FAIL, nil, c)
		return
	}

	//获取总数
	count, err := tagRepository.Count()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_COUNT_TAG_FAIL, nil, c)
		return
	}

	//组合数据
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = count
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

	//生成数据仓库
	tagRepository := repository.Tag{
		Name:      name,
		State:     1,
	}

	//判断标签是否存在
	exits := tagRepository.ExistByName()
	if exits {
		app.Response(http.StatusInternalServerError, errors. ERROR_EXIST_TAG, nil, c)
		return
	}

	//判断是否添加成功
	err := tagRepository.Add()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_ADD_TAG_FAIL, nil, c)
		return
	}

	//返回结果
	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
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
	tagRepository := repository.Tag{
		ID:id,
		Name:name,
	}
	tag := tagRepository.ExistById()
	if !tag {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	//组合数据
	err := tagRepository.Edit()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_EDIT_TAG_FAIL, nil, c)
		return
	}
	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
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
	tagRepository := repository.Tag{ID:id}
	tag := tagRepository.ExistById()
	if !tag {
		app.Response(http.StatusBadRequest, errors.ERROR_NOT_EXIST_TAG, nil, c)
		return
	}

	err := tagRepository.Delete()
	if err != nil {
		app.Response(http.StatusInternalServerError, errors.ERROR_DELETE_TAG_FAIL, nil, c)
		return
	}
	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
}

// @Tags 导出标签
// @Summary 导出文章标签
// @Param id formData string false "name"
// @Param state formData int false "state"
// @Param token path string true "token"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/export [Post]
func ExportTag(c *gin.Context)  {
	name := c.PostForm("name")
	state := 1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	//生成数据仓库
	tagRepository := repository.Tag{
		Name:      name,
		State:     state,
	}

	filename, err := tagRepository.Export()
	if err != nil {
		app.Response(http.StatusOK, errors.ERROR_EXPORT_TAG_FAIL, nil, c)
		return
	}

	app.Response(http.StatusOK, errors.SUCCESS, map[string]string{
		"export_path" : config.Apps.ExportPath + filename,
		"export_url" : app.GetExcelFullURL(filename),
	}, c)
}

// @Tags 导入标签
// @Summary 导入文章标签
// @Param file formData string false "file"
// @Param token path string true "token"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/import [Post]
func ImportTag(c *gin.Context)  {
	_, info, err := c.Request.FormFile("file")
	if err != nil {
		logger.Warn(err)
		app.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil, c)
		return
	}

	//判断文件大小
	//if !(info.Size >= config.Apps.UploadSize) {
	//	app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_MAX_SIZE_FAIL, nil, c)
	//	return
	//}
	logger.Info(info.Size)

	//判断后缀名
	ext := path.Ext(info.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_EXT_FAIL, nil, c)
		return
	}

	err = c.SaveUploadedFile(info, "public/" + info.Filename)
	if  err != nil {
		logger.Warn(err)
		app.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_SAVE_FAIL, nil, c)
		return
	}

	tagRepository := repository.Tag{}
	err = tagRepository.Import("public/" + info.Filename)
	if err != nil {
		logger.Warn(err)
		app.Response(http.StatusInternalServerError, errors.ERROR_IMPORT_TAG_FAIL, nil, c)
		return
	}

	app.Response(http.StatusOK, errors.SUCCESS, nil, c)
}