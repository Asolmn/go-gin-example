package tag_service

import (
	"encoding/json"
	"fmt"
	"github.com/Asolmn/go-gin-example/models"
	"github.com/Asolmn/go-gin-example/pkg/export"
	"github.com/Asolmn/go-gin-example/pkg/gredis"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/service/cache_service"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"time"
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

// 导出标签
func (t *Tag) Export() (string, error) {

	tags, err := t.GetAll() // 获取标签信息
	if err != nil {
		return "", err
	}
	file := xlsx.NewFile() // 创建excel文件
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow() // 添加行

	var cell *xlsx.Cell // 单元格
	for _, title := range titles {
		cell = row.AddCell() // 当前行添加一个单元格，并赋值给cell
		cell.Value = title   // 设置当前单元格的值
	}

	// 将tag信息添加到excel中
	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}
		row = sheet.AddRow() // 创建新一行
		// 将values的值逐一添加到新一行中的单元格中
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	times := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + times + ".xlsx" // 生成excel文件名

	fullPath := export.GetExcelFullPath() + filename // 设置完成路径
	err = file.Save(fullPath)                        // 保存
	if err != nil {
		return "", err
	}
	return filename, nil
}

// 导入标签
func (t *Tag) Import(r io.Reader) error {
	xlsxfile, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows, err := xlsxfile.GetRows("标签信息") // rows为sheet中的全部行
	if err != nil {
		return err
	}

	for irow, row := range rows {
		if irow > 0 { // 如果行数大于0
			var data []string
			for _, cell := range row { // 获取每一行的单元格值
				data = append(data, cell) // 将单元格的值循环加入到data中
			}
			// AddTag(名称,状态,创建人)
			err := models.AddTag(data[1], 1, data[2])
			if err != nil {
				return err
			}
		}
	}
	return nil
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
