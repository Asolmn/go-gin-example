package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()                                                               // 日志文件存放前缀
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt) // 日志文件命名形式

	// 例如：runtime/logs/log20230403.log
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath) // 返回文件信息结构描述文件

	switch {
	case os.IsNotExist(err): // 目录是否存在
		mkDir()
	case os.IsPermission(err): // 权限是否满足
		log.Fatalf("Permission :%v", err)
	}
	// 调用文件, 传入文件名称，指定的模式，文件权限
	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}
	return handle
}

func mkDir() {
	dir, _ := os.Getwd() // 返回当前目录对应的根路径名
	// os.ModePerm const定位ModePerm FileMode = 0777
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm) // 创建对应的目录以及所需的子目录
	if err != nil {
		panic(err)
	}
}
