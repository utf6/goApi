package app

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/pkg/config"
	errors "github.com/utf6/goApi/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

// hash 加密
func HashAndSalt(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {

	}
	return string(hash)
}

// hash 加密验证
func ValidatePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}

// md5 加密
func Md5(value string) string {
	str := md5.New()
	str.Write([]byte(value))

	return hex.EncodeToString(str.Sum(nil))
}

//返回结果
func Response(httpCode, code int, data interface{}, C *gin.Context) {
	C.JSON(httpCode, gin.H{
		"code" : code,
		"msg" : errors.GetMsg(code),
		"data": data,
	})
}

//获取excle 表格路由
func GetExcelFullURL(name string) string {
	return config.Apps.ImageUrl + "/" + config.Apps.ExportPath + name
}

func GetExcelFullPath() string {
	return config.Apps.RootPath + config.Apps.ExportPath
}