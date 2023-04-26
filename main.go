package main

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/models"
	"github.com/Asolmn/go-gin-example/pkg/gredis"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/routers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"syscall"
)

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	setting.Setup()
	models.Setup()
	logging.Setup()
	err := gredis.Setup()
	if err != nil {
		return
	}

	// 返回type Engine struct类型，包含RouterGroup
	// 相当于创建一个路由Handlers，后期可以绑定各类的路由规则和函数、中间件
	//routersInit := routers.InitRouter()
	//readTimeout := setting.ServerSetting.ReadTimeout
	//writeTimeout := setting.ServerSetting.WriteTimeout
	//endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	//maxHeaderBytes := 1 << 20
	//
	//server := &http.Server{
	//	Addr:           endPoint,
	//	Handler:        routersInit,
	//	ReadTimeout:    readTimeout,
	//	WriteTimeout:   writeTimeout,
	//	MaxHeaderBytes: maxHeaderBytes,
	//}
	//
	//log.Printf("[info] start http server listening %s", endPoint)
	//
	//server.ListenAndServe()

	// 通过endless实现服务重启的零停机
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	// NewServer返回一个初始化的endlessServer对象。在上面调用Serve实际上会“启动”服务器。
	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid()) // 输出进程的pid
	}

	// 与router.Run()方式本质上没有区别
	err = server.ListenAndServe()
	if err != nil {
		return
	}

}
