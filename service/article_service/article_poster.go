package article_service

import (
	"github.com/Asolmn/go-gin-example/pkg/file"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

// 合并图像标志
func GetPosterFlag() string {
	return "poster"
}

// 检查合并图像是否存在
func (a *ArticlePoster) CheckMergedImage(path string) bool {
	// 检查runtime/qrcode/PosterName是否存在
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}
	return true
}

// 打开合并图像
func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	// path = runtime/qrcode/
	// a.PosterName = poster-MD5(QRCODE_URL).jpg
	// MustOpen(runtime/qrcode/poster-MD5(QRCODE_URL).jpg)
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath() // runtime/qrcode/ 获取二维码存储路径
	// fileName = util.EncodeMD5(QrCode.URL).jpg
	// path = runtime/qrcode/
	// path+fileName = 新生成二维码的路径
	fileName, path, err := a.ArticlePoster.Qr.Encode(fullPath) // 生成二维码图像

	if err != nil {
		return "", "", err
	}

	if !a.CheckMergedImage(path) { // 检查合并后图像是否存在
		mergedF, err := a.OpenMergedImage(path) // 生成待合并图像mergedF
		if err != nil {
			return "", "", err
		}

		defer func(mergedF *os.File) { // 关闭mergedF
			err := mergedF.Close()
			if err != nil {
				logging.Info(err)
			}
		}(mergedF)

		bgF, err := file.MustOpen(a.Name, path) // 打开背景图
		if err != nil {
			return "", "", err
		}
		defer func(bgF *os.File) { // 关闭bgF
			err := bgF.Close()
			if err != nil {
				logging.Info(err)
			}
		}(bgF)

		qrF, err := file.MustOpen(fileName, path) // 打开生成的二维码图像
		if err != nil {
			return "", "", err
		}
		defer func(qrF *os.File) { // 关闭qrF
			err := qrF.Close()
			if err != nil {
				logging.Info(err)
			}
		}(qrF)

		bgImage, err := jpeg.Decode(bgF) // 解码bgF
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF) // 解码qrF
		if err != nil {
			return "", "", err
		}
		// 创建一个新的RGBA图像
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		// 在RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		err = jpeg.Encode(mergedF, jpg, nil)
		if err != nil {
			return "", "", err
		}
	}
	// 返回合并后图像的文件名和路径还有报错
	return fileName, path, nil

}
