package repository

import (
	"encoding/json"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/cache"
	"github.com/utf6/goApi/pkg/logger"
	"strconv"
	"strings"
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
