package article_service

import (
	"encoding/json"
	"fmt"
	"github.com/Asolmn/go-gin-example/models"
	"github.com/Asolmn/go-gin-example/pkg/gredis"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	CreatedBy     string
	ModifiedBy    string
	State         int

	PageNum  int
	PageSize int
}

// 添加文章
func (a *Article) Add() error {
	article := map[string]interface{}{ // 初始化文章结构体
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}
	if err := models.AddArticle(article); err != nil { // 调用数据库添加文章函数
		return err
	}
	return nil
}

// 编辑文章
func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	})
}

// 获取文章信息
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID} // 设置文章缓存
	key := cache.GetArticleKey()             // 获取文章缓存的key
	fmt.Println(key)
	if gredis.Exists(key) { // 判断文章缓存是否存在
		data, err := gredis.Get(key) // 获取文章缓存
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticle) // 对json进行反序列化，将json中的内容填到Article结构体中
			if err != nil {
				return nil, err
			}
			return cacheArticle, nil // 返回Article结构体
		}
	}
	article, err := models.GetArticle(a.ID) // 通过文章id获取数据
	if err != nil {
		return nil, err
	}
	err = gredis.Set(key, article, 300) // 设置缓存
	if err != nil {
		return nil, err
	}
	return article, nil
}

// 获取所有文章
func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey() // 返回所有文章缓存的key
	fmt.Println(key)
	if gredis.Exists(key) { // 判断key是否存在
		data, err := gredis.Get(key) // 获取对应的value
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticles) // 反序列化
			if err != nil {
				return nil, err
			}
			return cacheArticles, nil
		}
	}
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps()) // 查询全部文章
	if err != nil {
		return nil, err
	}
	err = gredis.Set(key, articles, 300) // 设置缓存
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// 删除文章
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

// 检测文章是否存在
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

// 获取文章数量
func (a *Article) Count() (int64, error) {
	return models.GetArticleTotal(a.getMaps())
}

// 结构体转换为map类型
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}
