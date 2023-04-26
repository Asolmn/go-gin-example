package tag_service

import (
	"encoding/json"
	"fmt"
	"github.com/Asolmn/go-gin-example/models"
	"github.com/Asolmn/go-gin-example/pkg/gredis"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

// 标签名是否存在
func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

// 标签id是否存在
func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

// 统计总共标签数量
func (t *Tag) Count() (int64, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag // 设置保存数据的变量
	)

	cache := cache_service.Tag{
		State: t.State,

		PageSize: t.PageSize,
		PageNum:  t.PageNum,
	} // 缓存数据

	if t.Name != "" {
		cache.Name = t.Name
	}

	key := cache.GetTagsKey() // 获取缓存的key
	fmt.Println(key)
	if gredis.Exists(key) { // 校验key
		data, err := gredis.Get(key) // 获取缓存数据
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheTags) // 反序列化,json转为结构体
			if err != nil {
				return nil, err
			}
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps()) // 获取数据
	if err != nil {
		return nil, err
	}
	err = gredis.Set(key, tags, 300) // 设置缓存
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// 创建标签
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

// 编辑标签
func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy // 修改用户
	data["name"] = t.Name
	if t.State >= 0 { // 检验状态
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

// 删除标签
func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

// 结构体序列化为map
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}
