package logger

import (
	"fmt"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/pkg/files"
	"os"
	"time"
)

//获取日志路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", config.Apps.RuntimePath, config.Apps.LogPath)
}

//获取日志全路径
func getLogFileName() string {
	return fmt.Sprintf("%s.%s", time.Now().Format(config.Apps.TimeFormat), config.Apps.LogExt)
}

//打开文件
func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := files.CheckPermission(src)

	if perm == true {
		return nil, fmt.Errorf("files.CheckPermission Permission denied src: %s", src)
	}

	err = files.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("files.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := files.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Faild to OpenFile :%v", err)
	}

	return f, nil
}
