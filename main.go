package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// http服务
var srv *http.Server

func main() {
	app := gin.Default()
	// Server Info
	app.GET("/", GetSendCount)
	app.GET("/send", SendMsg)
	srv = &http.Server{
		Addr:    "0.0.0.0:88",
		Handler: app,
	}

	//启动http请求
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	shutdown()
}

// 优雅的关闭
func shutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Errorf("Server Shutdown:[%v]", err)
		return
	}
}
