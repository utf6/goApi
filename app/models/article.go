package models

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Thumb string `json:"thumb"`
	Content string `json:"content"`
	State   int    `json:"state"`
}

//判断文章是否存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error

	if err != nil {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

//获取文章总数
func GetArticleTotal(maps interface{}) (int, error) {
	var  count int
	err := db.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

//获取所有文章
func GetArticles(pageNum int, pageSize int, maps interface{}) ([] *Article, error) {
	var articles [] *Article

	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&articles).Error
	if err != nil {
		return  nil, err
	}

	return articles, nil
}

//获取单个文章
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? and state = ?", id, 1).First(&article).Error

	if err != nil && article.ID > 0{
		 return nil, err
	}

	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil{
		return  nil, err
	}

	return &article, nil
}

//编辑文章
func EditArticle(id int, data interface{}) error {
	result := db.Model(&Article{}).Where("id = ?", id).Updates(data)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

//添加文章
func AddArticle(data map[string]interface{}) error {
	result := db.Create(&Article{
		TagID:   data["tag_id"].(int),
		Title:   data["title"].(string),
		Desc:    data["desc"].(string),
		Content: data["content"].(string),
		State:   1,
	})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

//删除文章
func DeleteArticle(id int) error {
	result := db.Where("id = ?", id).Delete(&Article{})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CleanArticle() (bool, error) {
	result := db.Unscoped().Where("state = ?", -1).Delete(&Article{})
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
