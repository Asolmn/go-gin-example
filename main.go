package main

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/models"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {

	setting.Setup()
	fmt.Println(setting.AppSetting.LogSaveName)
	models.Setup()
	logging.Setup()

	// 通过endless实现服务重启的零停机
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	// 返回type Engine struct类型，包含RouterGroup
	// 相当于创建一个路由Handlers，后期可以绑定各类的路由规则和函数、中间件
	//router := routers.InitRouter()

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid()) // 输出进程的pid
	}
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HTTPPort), // 监听的tcp地址，格式为:8000
	//	Handler:        router,                               // http句柄，实质为ServerHTTP，用于处理程序响应http请求
	//	ReadTimeout:    setting.ReadTimeout,                  // 允许读取的最大时间
	//	WriteTimeout:   setting.WriteTimeout,                 // 允许写入的最大时间
	//	MaxHeaderBytes: 1 << 20,                              // 请求头的最大字节数 1*2^20=1MB
	//}

	// 与router.Run()方式本质上没有区别
	err := server.ListenAndServe()
	if err != nil {
		return
	}

}
