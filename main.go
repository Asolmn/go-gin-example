package main

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/routers"
	"net/http"
)

func main() {

	// 返回type Engine struct类型，包含RouterGroup
	// 相当于创建一个路由Handlers，后期可以绑定各类的路由规则和函数、中间件
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort), // 监听的tcp地址，格式为:8000
		Handler:        router,                               // http句柄，实质为ServerHTTP，用于处理程序响应http请求
		ReadTimeout:    setting.ReadTimeout,                  // 允许读取的最大时间
		WriteTimeout:   setting.WriteTimeout,                 // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20,                              // 请求头的最大字节数
	}

	// 与router.Run()方式本质上没有区别
	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
