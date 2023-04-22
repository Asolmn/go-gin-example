package api

import (
	"github.com/Asolmn/go-gin-example/pkg/e"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image") // 通过表单获取文件
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename) // 图片文件名
		fullPath := upload.GetImageFullPath()            // 图片完整url runtime/upload/images/
		savePath := upload.GetImagePath()                // 图片保存路径 upload/images/

		src := fullPath + imageName // runtime/upload/images/imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			// 检查图片后缀和图片大小
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath) // 检查图片存放目录
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil { // SaveUploadedFile将文件保存到指定src中
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName) // http://127.0.0.1:8000/upload/images/imageName
				data["image_save_url"] = savePath + imageName         // upload/images/imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
