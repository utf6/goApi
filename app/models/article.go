package models

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	State   int    `json:"state"`
}

//判断文章是否存在
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}
	return false
}

//获取文章总数
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

//获取所有文章
func GetArticles(pageNum int, pageSize int, maps interface{}) (article []Article) {
	db.Preload("Tag.name").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&article)
	return
}

//获取单个文章
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)

	return
}

//编辑文章
func EditArticle(id int, data interface{}) bool {
	result := db.Model(&Article{}).Where("id = ?", id).Updates(data)

	if result.Error != nil {
		return false
	}
	return true
}

//添加文章
func AddArticle(data map[string]interface{}) bool {
	result := db.Create(&Article{
		TagID:   data["tag_id"].(int),
		Title:   data["title"].(string),
		Desc:    data["desc"].(string),
		Content: data["content"].(string),
		State:   1,
	})

	if result.Error != nil {
		return false
	}
	return true
}

//删除文章
func DeleteArticle(id int) bool {
	result := db.Where("id = ?", id).Delete(&Article{})

	if result.Error != nil {
		return false
	}
	return true
}

func CleanArticle() bool {
	result := db.Unscoped().Where("state = ?", -1).Delete(&Article{})
	if result.Error != nil {
		return false
	}
	return true
}
