package models

type Tag struct {
	Model

	Name      string `json:"name"`
	State     int    `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) (tag []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tag)
	return
}

//获取标签总数
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
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
func AddTag(name string) bool {
	result := db.Create(&Tag{
		Name:      name,
		State:     1,
	})

	if result.Error != nil {
		return false
	}
	return true
}

//编辑标签
func EditTag(id int, data interface{}) bool {
	result := db.Model(&Tag{}).Where("id = ?", id).Updates(data)

	if result .Error != nil {
		return  false
	}

	return true
}

//删除标签
func DeleteTag(id int) bool {
	result := db.Where("id = ?", id).Delete(&Tag{})

	if result.Error != nil {
		return  false
	}
	return  true
}