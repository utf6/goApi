package config

import (
	"github.com/go-ini/ini"
	"github.com/utf6/goApi/pkg/logger"
	"time"
)

type App struct {
	JwtSecret string
	PageSize int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath string
	ImageMaxSize int
	ImageAllows []string

	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string
}

var Apps = &App{}

type Server struct {
	RunMode string
	HttpPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

var Servers = &Server{}

type Database struct {
	Type string
	User string
	Password string
	Host string
	Name string
	TablePrefix string
}

var Databases = &Database{}

func Setup() {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		logger.Error("Fail to parse 'conf/app.ini': %v", err)
	}

	//加载系统配置
	err = Cfg.Section("app").MapTo(Apps)
	if err != nil {
		logger.Error("Cfg.MapTo Apps err: %v", err)
	}

	//加载服务配置
	err = Cfg.Section("server").MapTo(Servers)
	if err != nil {
		logger.Error("Cfg.MapTo Servers err: %v", err)
	}

	//加载数据库配置
	err = Cfg.Section("database").MapTo(Databases)
	if err != nil {
		logger.Error("Cfg.MapTo Databasess err: %v", err)
	}
}
