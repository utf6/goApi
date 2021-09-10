package files

import (
	"fmt"
	"github.com/utf6/goApi/app"
	"github.com/utf6/goApi/pkg/config"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

//获取图片完整访问 URL
func GetImageFullUrl(name string) string {
	return config.Apps.ImageUrl + "/" + GetImagePath() + name
}

//获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := app.Md5(strings.TrimSuffix(name, ext))

	return fileName + ext
}

//获取图片路径
func GetImagePath() string {
	return config.Apps.ImageSavePath
}

func CheckImageExt(fileName string) bool {
	ext := GetExt(fileName)
	for _, allowExt := range config.Apps.ImageAllows {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := GetSize(f)
	if err != nil {
		log.Println(err)
		return false
	}
	return size <= config.Apps.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("files.IsNotExistMkDir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("files.CheckPermission Permission denied src: %s", src)
	}

	return nil
}