package models

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Model
	TagID int `json:"tag_id" gorm:"index"` // grom:index用于声明索引
	Tag   Tag `json:"tag"`                 // 嵌套Tag struct，利用TagID与Tag模型相互关联，执行查询的时候，能达到Article和Tag关联查询

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 检测文章是否存在
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}
	return false
}

// 统计文章个数
func GetArticleTotal(maps interface{}) (count int64) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

// 获取多个文章
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

// 获取单个文章
func GetArticle(id int) (article Article) {
	// 根据id查询文章
	db.Where("id = ?", id).First(&article)
	// Preload预加载Tag模型，这样可以在查询Article模型时同时查询关联的Tag模型
	db.Preload("Tag").Find(&article.Tag, article.TagID)
	return
}

// 编辑文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return true
}

// 添加文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})
	return true
}

// 删除文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}

func (article *Article) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (article *Article) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
