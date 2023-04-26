// 工具包-获取分页页码

package util

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int() // 获取的页码从字符串转成int型
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	fmt.Println(result)
	return result
}
