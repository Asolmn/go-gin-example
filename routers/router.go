package routers

import (
	_ "github.com/Asolmn/go-gin-example/docs"
	"github.com/Asolmn/go-gin-example/middleware/jwt"
	"github.com/Asolmn/go-gin-example/pkg/export"
	"github.com/Asolmn/go-gin-example/pkg/qrcode"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/pkg/upload"
	"github.com/Asolmn/go-gin-example/routers/api"
	v1 "github.com/Asolmn/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.ServerSetting.RunMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 创建不同的HTTP方法绑定到Handlers中
	// gin.H{}是一个map[string]interface{}类型
	// gin.Context是gin中的上下文，允许在中间件之间传递变量、管理流、验证json请求，响应等
	//r.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	// 创建静态文件服务器
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/auth", api.GetAuth)
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT()) // 将中间件接入Gin

	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新建表
		apiv1.POST("/tags", v1.AddTag)
		// 更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		// 导出标签
		apiv1.POST("/tags/export", v1.ExportTag)
		// 导入标签
		apiv1.POST("/tags/import", v1.ImporTag)
	}

	{
		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		// 生成二维码
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
