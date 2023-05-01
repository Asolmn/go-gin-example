package qrcode

import (
	"github.com/Asolmn/go-gin-example/pkg/file"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
	"os"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// 获取二维码存储相对地址
func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath // qrcode/
}

// 获取二维码存储完整地址
func GetQrCodeFullPath() string {
	// runtime/qrcode/
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

// 获取二维码完成URL地址
func GetQrCodeFullUrl(name string) string {
	// http://127.0.0.1:8000/qrcode/name
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

// 生成二维码文件名
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

// 二维码文件后缀
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// 检查文件是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}
	return true
}

// 生成二维码
func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()

	// runtime/qrcode/filename.jpg
	src := path + name // 获取二维码生成路径

	if file.CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode) // 创建二维码
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height) // 缩放二维码到指定大小
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path) // 新建存放二维码图片的文件
		if err != nil {
			return "", "", err
		}
		defer func(f *os.File) { // 最后运行关闭文件
			err := f.Close()
			if err != nil {

			}
		}(f)

		err = jpeg.Encode(f, code, nil) // 以JPEG 4：2：0基线格式写入文件
		if err != nil {
			return "", "", err
		}
	}
	// 返回二维码文件名，保存路径，报错
	return name, path, nil
}
