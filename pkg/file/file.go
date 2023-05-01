package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// 返回文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// 返回文件后缀
func GetExt(filename string) string {
	return path.Ext(filename)
}

// 返回文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// 创建文件
func MkDir(src string) error {
	// os.ModePerm const定位ModePerm FileMode = 0777
	// 创建对应的目录以及所需的子目录
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 根据特定模式打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 打开文件
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd() // 获取根路径
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath  // 绝对路径
	perm := CheckPermission(src) // 检查权限
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src) // 如果路径不存在则创建
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	// 打开/root/progject/go-gin-example/runtime/qrcode/fileName，如果文件不存在，则创建，权限为0644
	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile:%v", err)
	}
	return f, nil
}
