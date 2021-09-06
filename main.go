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

/**
package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/utf6/goApi/pkg/config"
	"github.com/utf6/goApi/routes"
	"log"
	"syscall"
)

func main() {
	endless.DefaultReadTimeOut = config.ReadTimeOut
	endless.DefaultWriteTimeOut = config.WriteTimeOut
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", config.HttpPort)

	server := endless.NewServer(endPoint, routes.InitRoute())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}


	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
*/
