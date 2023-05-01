package v1

import (
	"github.com/Asolmn/go-gin-example/pkg/app"
	"github.com/Asolmn/go-gin-example/pkg/e"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/qrcode"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/Asolmn/go-gin-example/pkg/util"
	"github.com/Asolmn/go-gin-example/service/article_service"
	"github.com/Asolmn/go-gin-example/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary 获取单个文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	// 限制ID格式
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	// 检测ID格式
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)                        // 调用生成报错函数
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil) // 返回响应
		return
	}

	articleService := article_service.Article{ID: id} // 生成servicec层的article结构体
	exists, err := articleService.ExistByID()
	if err != nil { // 检测存在文章失败
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists { // 文章不存在
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	article, err := articleService.Get() // 获取该id的文章数据
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article) // 返回响应
}

// @Summary 获取多个文章
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" { // 状态校验
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" { // 标签ID校验
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() { // 参数验证
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ // 生成servicec层的article结构体
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count() // 文章总数
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}
	articles, err := articleService.GetAll() // 全部文章数据
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["totail"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 添加文章表单
type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 新增文章
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)

	httpCode, errCode := app.BindAndValid(c, &form) // 根据表单参数返回http状态码和错误码

	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil { // 检查已存在标签失败
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists { // 标签不存在
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 更新文章表单
type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 更新文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form EditArticleForm
	)
	httpCode, errCode := app.BindAndValid(c, &form) // 根据表单参数返回http状态码和错误码

	if errCode != e.SUCCESS { // 验证表单参数
		appG.Response(httpCode, errCode, nil)
		return
	}
	articleService := article_service.Article{ // service的article结构体
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil { // 检验获取已存在文章是
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists { // 检查文章是否存在
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()

	if err != nil { // 检验获取已存在标签是
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists { // 检查标签是否存在
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit() // 更新文章
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0") // 校验参数

	if valid.HasErrors() { // 对请求参数处理
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil { // 检查文章是否存在
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists { // 文章不存在
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	err = articleService.Delete()
	if err != nil { // 检查删除情况
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

type qrcodeUrl struct {
	QRCODE_URl string
}

func GenerateArticlePoster(c *gin.Context) {
	var QRCODE_URL string = "http://github.com/Asolmn/go-gin-example"

	appG := app.Gin{C: c}

	article := &article_service.Article{}
	qrC := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto) // 创建二维码信息

	// posterName = poster-MD5(QRCODE_URL).jpg
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qrC.URL) + qrC.GetQrCodeExt()

	/*
		articlePoster = {
			PosterName string = poster-utile.EncodeMD5(QRCODE_URL).jpg
			Article *Atricle = Article{}
			qr *qrcode.QrCode = qrC {
					URL    string = QRCODE_URL
					Width  int = 300
					Height int = 300
					Ext    string = .jpg
					Level  qr.ErrorCorrectionLevel = qr.M
					Mode   qr.Encoding = qr.Auto
			}
		}
	*/
	articlePoster := article_service.NewArticlePoster(posterName, article, qrC)
	/*
		articlePosterBgService = {
			Name string =  bg.jpg
			*ArticlePoster = articlePoster
			Rect: Rect{
				X0: 0,
				Y0: 0,
				X1: 550,
				Y1: 700,
			}
			Pt: Pt{
				X: 125,
				Y: 298,
			}
		}
	*/
	articlePosterBgService := article_service.NewArticlePosterBg( // 合并图像结构体
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate() // 生成背景与二维码的合并图像
	if err != nil {                                       // 返回报错
		logging.Info(err)
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POST_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{ // 返回成功响应
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
