package repository

import (
	"encoding/json"
	"fmt"
	"github.com/utf6/goApi/app"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/cache"
	"github.com/utf6/goApi/pkg/logger"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	ID int
	Name string
	State int

	CreatedAt string
	UpdatedAt string
	PageNum int
	PageSize int
}

//通过id判断标签是否存在
func (t *Tag) ExistById() bool {
	return models.ExistTagById(t.ID)
}

//获取所有标签
func (t *Tag) GetAll() ([]*models.Tag, error) {
	var (
		tags, tagsCache []*models.Tag
	)

	keys := t.GetCacheKeys()
	if cache.Exists(keys) {
		data, err := cache.Get(keys)
		if err != nil {
			logger.Error(err)
		} else {
			json.Unmarshal(data, &tagsCache)
			return tagsCache, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.GetMaps())
	if err != nil {
		return nil, err
	}

	cache.Set(keys, tags, 3600)
	return tags, nil
}

//获取总数
func (t *Tag) Count() (int ,error) {
	return models.GetTagTotal(t.GetMaps())
}

//通过名称判断标签是否存在
func (t *Tag) ExistByName() bool {
	return models.ExistTagByName(t.Name)
}

//添加标签
func (t *Tag) Add() error {
	return models.AddTag(t.Name)
}

//编辑标签
func (t *Tag) Edit() error {
	tag := models.Tag{
		Name:  t.Name,
	}
	return models.EditTag(t.ID, tag)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}


//组合条件
func (t *Tag) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}

	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

//获取缓存key
func (t *Tag) GetCacheKeys() string {
	keys := []string{"tags", "list"}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}

//导出标签
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := excelize.NewFile()
	index := file.NewSheet("sheet")

	//写入表头
	titles := map[string]string{"A1" : "ID", "B1" : "名称", "C1" : "创建时间", "D1" : "更新时间", "E1" : "状态"}
	for k, title := range titles {
		file.SetCellValue("sheet", k, title)
	}

	for id, tag := range tags {
		//组合数据
		staStr := "正常"
		if tag.State == 0 {
			staStr = "禁用"
		} else if tag.State == -1 {
			staStr = "删除"
		}

		values := map[string]string{
			fmt.Sprintf("A%d", id+1) : strconv.Itoa(tag.ID),
			fmt.Sprintf("B%d", id+1) : tag.Name,
			fmt.Sprintf("C%d", id+1) : tag.CreatedAt.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("D%d", id+1) : tag.UpdatedAt.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("E%d", id+1) :  staStr, //strconv.Itoa(tag.State),
		}

		//写入表格
		for ck, value := range values {
			file.SetCellValue("sheet", ck, value)
		}
	}

	file.SetActiveSheet(index)

	saveName := "tags_" + strconv.Itoa(int(time.Now().Unix())) + ".xlsx"
	savePath := app.GetExcelFullPath() + saveName  //app.GetExcelFullPath()

	// Save spreadsheet by the given path.
	if err := file.SaveAs(savePath); err != nil {
		return "", err
	}
	return saveName, nil
}