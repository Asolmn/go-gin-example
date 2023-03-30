package routers

import (
	"github.com/Asolmn/go-gin-example/pkg/setting"
	v1 "github.com/Asolmn/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.RunMode)
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

	apiv1 := r.Group("/api/v1")

	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新建表
		apiv1.POST("/tags", v1.AddTag)
		// 更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
	}

	return r
}