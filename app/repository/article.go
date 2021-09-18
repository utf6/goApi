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

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	Thumb string
	State         int
	CreatedAt     string
	UpdateAt    string

	PageNum  int
	PageSize int
}

//新增
func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id": a.TagID,
		"title": a.Title,
		"desc" : a.Desc,
		"content": a.Content,
		"thumb": a.Thumb,
		"state": a.State,
		"created_at" : a.CreatedAt,
		"updated_at": a.UpdateAt,
	}
	err := models.AddArticle(article);
	if err != nil {
		return err
	}
	return  nil
}

//编辑
func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id": a.TagID,
		"title": a.Title,
		"desc" : a.Desc,
		"content": a.Content,
		"thumb": a.Thumb,
		"state": a.State,
		"updated_at": a.UpdateAt,
	})
}

//获取单条数据
func (a *Article) Get() (*models.Article, error) {
	var articleCache *models.Article

	key := "article_" + strconv.Itoa(a.ID);
	if cache.Exists(key) {
		data, err := cache.Get(key)
		if err != nil {
			logger.Info(err)
		} else  {
			json.Unmarshal(data, &articleCache)
			return  articleCache, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return  nil, err
	}
	cache.Set(key, article, 3600)
	return article, nil
}

//获取所有文章
func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, articlesCache []*models.Article
	)

	//获取缓存key
	keys := a.GetCacheKeys()
	//如果缓存存在，直接返回缓存
	if cache.Exists(keys) {
		data, err := cache.Get(keys)
		if err != nil {
			logger.Error(err)
		} else  {
			json.Unmarshal(data, &articlesCache)
			return articlesCache, nil
		}
	}

	//缓存不存在则读取数据
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.GetMaps())
	if err != nil {
		return nil, err
	}

	cache.Set(keys, articles, 3600)
	return articles, nil
}

//获取文章数量
func (a *Article) Count() (int, error) {
	var count int
	count, err := models.GetArticleTotal(a.GetMaps())
	if err != nil {
		return 0, err
	}

	return count, nil
}

//删除数据
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

//判断文章是否存在
func (a *Article) ExistByID() bool {
	return models.ExistArticleByID(a.ID)
}

//组合条件
func (a *Article) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.State >= 0 {
		maps["state"] = a.State
	}
	if a.TagID > 0 {
		maps["tag_id"] = a.TagID
	}
	return maps
}

//获取缓存key
func (a *Article) GetCacheKeys() string {
	keys := []string{"article", "list"}

	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}

//导入文章
func (a *Article) Import(filename string) error {

	var file, err = excelize.OpenFile(filename)
	if err != nil {
		return  err
	}

	var rows, _ = file.GetRows("Sheet1")
	for key, row := range rows {
		if key > 0 {
			//组合数据
			TagID, _ := strconv.Atoi(row[3])
			var article = map[string]interface{}{
				"tag_id":  TagID,
				"title":   row[0],
				"desc":    row[1],
				"content": row[2],
				"state":   1,
			}
			_ = models.AddArticle(article)
		}
	}
	return  nil
}

//导出文章
func (a *Article) Export() (string, error) {
	articles, err := a.GetAll()
	if err != nil {
		return "", err
	}

	file := excelize.NewFile()
	index := file.NewSheet("Sheet1")

	//写入表头
	titles := map[string]string{"A1" : "ID", "B1" : "标题", "C1" : "描述", "D1" : "内容", "E1" : "状态", "F1" : "标签id", "G1" : "缩略图", "H1" : "创建时间", "I1" : "更新时间"}
	for k, title := range titles {
		file.SetCellValue("Sheet1", k, title)
	}

	for id, article := range articles {
		//组合数据
		staStr := "正常"
		if article.State == 0 {
			staStr = "禁用"
		} else if article.State == -1 {
			staStr = "删除"
		}

		var (
			values = map[string]string{
				fmt.Sprintf("A%d", id+2): strconv.Itoa(article.ID),
				fmt.Sprintf("B%d", id+2): article.Title,
				fmt.Sprintf("C%d", id+2): article.Desc,
				fmt.Sprintf("D%d", id+2): article.Content,
				fmt.Sprintf("E%d", id+2): staStr,//strconv.Itoa(tag.State),
				fmt.Sprintf("F%d", id+2): strconv.Itoa(article.TagID),
				fmt.Sprintf("G%d", id+2): article.Thumb,
				fmt.Sprintf("H%d", id+2): article.CreatedAt.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("I%d", id+2): article.UpdatedAt.Format("2006-01-02 15:04:05"),
			} //写入表格
		)
		for ck, value := range values {
			file.SetCellValue("Sheet1", ck, value)
		}
	}

	file.SetActiveSheet(index)

	saveName := "articles_" + strconv.Itoa(int(time.Now().Unix())) + ".xlsx"
	savePath := app.GetExcelFullPath() + saveName  //app.GetExcelFullPath()

	// Save spreadsheet by the given path.
	if err := file.SaveAs(savePath); err != nil {
		return "", err
	}
	return saveName, nil
}