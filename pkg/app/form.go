package app

import (
	"github.com/Asolmn/go-gin-example/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form) // 根据内容类型选择绑定类型
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{} // 上下文管理数据验证和错误消息
	check, err := valid.Valid(form)  // 有效验证结构

	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check { // 如果校验不通过
		MarkErrors(valid.Errors) // 生成日志报错
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
