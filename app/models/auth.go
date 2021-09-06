package models

import "github.com/utf6/goApi/app"

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("id,password").Where("username = ?", username).First(&auth)

	if auth.ID < 0 && !app.ValidatePasswords(auth.Password, password) {
		return false
	}
	return true
}
