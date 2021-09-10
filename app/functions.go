package app

import (
	"crypto/md5"
	"encoding/hex"
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
