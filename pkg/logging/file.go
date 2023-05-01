package logging

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/file"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"os"
	"time"
)

// 获取日志文件路径
func getLogFilePath() string {
	// runtime/logs/
	return fmt.Sprintf("%s%s",
		setting.AppSetting.RuntimeRootPath,
		setting.AppSetting.LogSavePath)
}

// 获取日志文件名
func getLogFileName() string {
	// log2023430.log
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)
}

// 打开日志文件
func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	// src = /root/project/go-gin-example/runtime/logs/
	// 因为logging的初始化调用是在main.go中运行，所以通过os.Getwd()返回的当前目录相对应的根路径名为main.go文件所对应的根路径
	// 也就是/root/project/go-gin-example
	src := dir + "/" + filePath
	perm := file.CheckPermission(src)

	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src) // 如果路径不存在则进行创建
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	// open(/root/project/go-gin-example/runtime/logs/filename)
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}
