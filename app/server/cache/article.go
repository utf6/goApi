package cache

import (
	"strconv"
	"strings"
)

type Article struct {
	ID int
	TagID int
	State int
	PageNum int
	PageSize int
}

func (a *Article) GetKey() string {
	return "article_" + strconv.Itoa(a.ID)
}

func (a *Article) GetKeys() string {
	keys := []string{"article", "list"}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State > 0 {
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