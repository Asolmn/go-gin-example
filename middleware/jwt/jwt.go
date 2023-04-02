package jwt

import (
	"github.com/Asolmn/go-gin-example/pkg/e"
	"github.com/Asolmn/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token") // 获取token
		if token == "" {          // 如果token为空
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token) // 验证token
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL // token鉴权失败
			} else if time.Now().Unix() > claims.ExpiresAt { // 检测时间是否超过3个小时
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT // token超时
			}
		}

		if code != e.SUCCESS { // 返回验证失败信息
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort() // 终止请求处理
			return
		}
		c.Next() // 调用下一个中间件
	}
}
