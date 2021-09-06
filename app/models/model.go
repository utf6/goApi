package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/pkg/logger"
	"time"
)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_ke;auto" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State int `json:"state"`
}

func init() {
	sec, err := config.Cfg.GetSection("database")
	if err != nil {
		logger.Error("Fail to get section 'database': %v", err)
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
		logger.Error(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return sec.Key("TABLE_PREFIX").String() + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//注册Callbacks
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback) //删除数据时回调
}

func CloseDB() {
	defer db.Close()
}

//增加数据回调
func updateTimeStampForCreateCallback(scope *gorm.Scope)  {
	if !scope.HasError() {
		nowTime := time.Now()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(nowTime)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(nowTime)
			}
		}
	}
}

//更新数据回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("state")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}