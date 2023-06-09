package models

import (
	"gorm.io/gorm"
)

// 用于Gorm的使用 给予附属属性json，方便c.JSON的时候自动转换格式
type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 判断tag是否存在
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	// select id from tag where name = name order by blog_tags.id limit 1
	err := db.Select("id").Where("name = ? and deleted_on = ?", name, 0).First(&tag).Error
	//fmt.Println(tag.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 获取tag列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	if pageSize > 0 && pageNum > 0 {
		// 以maps为条件，offset指定开始返回记录前跳过的记录数
		// Limit表示语句中的limit 1
		// 查询结果放入tags中
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return
}

// 获取tag数量
func GetTagTotal(maps interface{}) (count int64, err error) {
	// 查询整个tag表，以maps为条件，返回符合条件的记录数量
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return
}

// 添加tag
func AddTag(name string, state int, createdBy string) error {
	t := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	result := db.Create(&t)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

// 检测id
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	// 查询id
	// select id from blog_tag where blog_tag.id = id and deleted_on = 0 order by id limit 1
	err := db.Select("id").Where("id = ? and deleted_on = ?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 删除tag
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}

	return nil
}

// 修改tag
func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// 硬删除
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})
	return true
}
