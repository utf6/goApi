package models

type Tag struct {
	Model

	Name  string `json:"name"`
	State int    `json:"state"`
}

//获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) ([] *Tag, error) {
	var tags [] *Tag
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Order("id DESC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

//获取标签总数
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

//通过名称判断标签是否存在
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)

	if tag.ID > 0 {
		return true
	}
	return false
}

//通过id判断标签是否存在
func ExistTagById(id int) bool {
	var tag Tag
	db.Where("id =  ?", id).First(&tag)

	if tag.ID > 0 {
		return true
	}
	return false
}

//新增标签
func AddTag(name string) error {
	err := db.Create(&Tag{
		Name:  name,
		State: 1,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

//编辑标签
func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error

	if err != nil {
		return err
	}

	return nil
}

//删除标签
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error

	if err != nil {
		return err
	}
	return nil
}

func CleanTag() bool {
	result := db.Unscoped().Where("state = ?", -1).Delete(&Tag{})
	if result != nil {
		return false
	}
	return true
}
