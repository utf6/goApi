package config

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	JwtSecret   string
	PageSize    int
	RuntimePath string
	ExportPath string
	RootPath string

	ImageUrl      string
	ImageSavePath string
	ImageMaxSize  int
	UploadSize int64
	ImageAllows   []string

	LogPath    string
	LogExt     string
	TimeFormat string
}

var Apps = &App{}

type Server struct {
	HttpPort     int
	RunMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var Servers = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var Databases = &Database{}


type GRedis struct {
	Host        string
	MaxIdle        int
	Password    string
	MaxActive        int
	IdleTimeout        time.Duration
}

var Redis = &GRedis{}

func Setup() {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		log.Printf("Fail to parse 'conf/app.ini': %v", err)
	}

	//加载系统配置
	err = Cfg.Section("app").MapTo(Apps)
	if err != nil {
		log.Printf("Cfg.MapTo Apps err: %v", err)
	}

	//加载数据库配置
	err = Cfg.Section("database").MapTo(Databases)
	if err != nil {
		log.Printf("Cfg.MapTo Databases err: %v", err)
	}

	//加载服务配置
	err = Cfg.Section("server").MapTo(Servers)
	if err != nil {
		log.Printf("Cfg.MapTo Servers err: %v", err)
	}

	//加载redis配置
	err = Cfg.Section("redis").MapTo(Redis)
	if err != nil {
		log.Printf("Cfg.MapTo Redis err: %v", err)
	}
}
