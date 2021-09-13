package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/utf6/goApi/app/models"
	"github.com/utf6/goApi/pkg/cache"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/pkg/logger"
	"github.com/utf6/goApi/routes"
	"net/http"
	"time"
)

func init() {
	config.Setup()
	models.Setup()
	logger.Setup()
	cache.Setup()
}

func main() {
	gin.SetMode(config.Servers.RunMode)

	r := routes.InitRoute()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Servers.HttpPort),
		Handler:        r,
		ReadTimeout:    config.Servers.ReadTimeout * time.Second,
		WriteTimeout:   config.Servers.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
