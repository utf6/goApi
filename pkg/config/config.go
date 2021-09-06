package config

import (
	"github.com/go-ini/ini"
	"github.com/utf6/goApi/pkg/logger"
	"time"
)

var (
	Cfg          *ini.File
	RunMode      string
	HttpPort     int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration

	PageSize  int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("app.ini")
	if err != nil {
		logger.Error("Fail to parse 'conf/app.ini': %v", err)
	}

	//加载运行模式
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	//加载服务配置
	ser, err := Cfg.GetSection("server")
	if err != nil {
		logger.Error("Fail to get section 'server': %v", err)
	}
	HttpPort = ser.Key("HTTP_PORT").MustInt(8000)
	ReadTimeOut = time.Duration(ser.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeOut = time.Duration(ser.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second

	//加载系统配置
	app, err := Cfg.GetSection("app")
	if err != nil {
		logger.Error("Fail to get section 'app': %v", err)
	}
	JwtSecret = app.Key("JWT_SECRET").MustString("saf2sa23@!##")
	PageSize = app.Key("PAGE_SIZE").MustInt(10)
}
