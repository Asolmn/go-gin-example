package models

import (
	"gorm.io/gorm"
	"log"
	"time"
)

// 用于Gorm的使用 给予附属属性json，方便c.JSON的时候自动转换格式
type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int64  `json:"state"`
}

// 获取tag列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	// 以maps为条件，offset指定开始返回记录前跳过的记录数
	// Limit表示语句中的limit 1
	// 查询结果放入tags中
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

// 获取tag数量
func GetTagTotal(maps interface{}) (count int64) {
	// 查询整个tag表，以maps为条件，返回符合条件的记录数量
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

// 判断tag是否存在
func ExistTagByName(name string) bool {
	var tag Tag
	// select id from tag where name = name order by blog_tags.id limit 1
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// 添加tag
func AddTag(name string, state int64, createdBy string) bool {

	t := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	result := db.Create(&t)
	if err := result.Error; err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// 检测id
func ExistTagByID(id int) bool {
	var tag Tag
	// 查询id
	// select id from blog_tag where blog_tag.id = id order by id limit 1
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// 删除tag
func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

// 修改tag
func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)
	return true
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {

	tx.Statement.SetColumn("CreatedOn", time.Now().Unix())
	return
}

func (tag *Tag) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("ModifiedOn", time.Now().Unix())
	return
}
