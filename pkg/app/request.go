package app

import (
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/astaxie/beego/validation"
)

// MarkErrors记录错误日志
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
