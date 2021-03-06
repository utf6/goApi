package routes

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/utf6/goApi/app"
	"github.com/utf6/goApi/app/http/controllers/api"
	v1 "github.com/utf6/goApi/app/http/controllers/api/v1"
	"github.com/utf6/goApi/app/http/middleware"
	_ "github.com/utf6/goApi/docs"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/pkg/files"
	"net/http"
)

func InitRoute() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.Servers.RunMode)

	r.POST("/auth/getToken", api.GetToken)
	r.POST("/auth/uploads", api.Uploads)
	r.StaticFS("/uploads/images", http.Dir(config.Apps.RootPath + files.GetImagePath()))
	r.StaticFS("/export", http.Dir(app.GetExcelFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.Auth())
	{
		//获取标签列表
		apiV1.GET("/tags", v1.GetTags)
		//新建标签
		apiV1.POST("/tags", v1.AddTag)
		//导出标签
		apiV1.GET("/tags/export", v1.ExportTag)
		//导入标签
		apiV1.POST("/tags/import", v1.ImportTag)
		//更新指定标签
		apiV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiV1.POST("/articles", v1.AddArticle)
		//导入文章
		apiV1.POST("/articles/import", v1.ImportArticle)
		//文章
		apiV1.GET("/articles/export", v1.ExportArticle)
		//更新指定文章
		apiV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}
	return r
}
