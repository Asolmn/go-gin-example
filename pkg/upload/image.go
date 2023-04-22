package upload

import (
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/file"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 获取图片完整访问URL
func GetImageFullUrl(name string) string {
	// http://127.0.0.1:8000/upload/images/name
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

// 获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName) // 使用md5对图片名称进行转换

	// md5(filename).jpg,.png,.jpeg
	return fileName + ext
}

// 获取图片路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath // upload/images/
}

// 获取图片完整路径
func GetImageFullPath() string {
	// runtime/upload/images/
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// 检查图片后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName) // 获取后缀
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

// 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

// 检查图片存放目录
func CheckImage(src string) error {
	dir, err := os.Getwd() // /root/project/go-gin-example
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	// /root/project/go-gin-example/runtime/upload/images
	err = file.IsNotExistMkDir(dir + "/" + src) // 如果不存在，则创建新目录和子目录
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src) // 检查权限
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
