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

//获取多个文章标签
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

//新增文章标签
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
			} else  {
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

//修改文章标签
func EditTag(c *gin.Context) {
}

//删除文章标签
func DeleteTag(c *gin.Context) {
}
