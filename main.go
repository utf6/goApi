package main

import (
	"fmt"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/routes"
	"net/http"
)

func main() {
	r := routes.InitRoute()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.HttpPort),
		Handler:        r,
		ReadTimeout:    config.ReadTimeOut,
		WriteTimeout:   config.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
