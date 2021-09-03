package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/utf6/goApi/pkg/config"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID        int `gorm:"primary_ke;auto" json:"id"`
}

func init() {
	sec, err := config.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	db, err = gorm.Open(
		sec.Key("TYPE").String(),
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			sec.Key("USER").String(),
			sec.Key("PASSWORD").String(),
			sec.Key("HOST").String(),
			sec.Key("NAME").String(),
		))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return sec.Key("TABLE_PREFIX").String() + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now().Format("2006-01-02 15:04:05"))
	scope.SetColumn("UpdatedAt", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}