package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MServer interface {
	Start()
	GracefulStop()
}

var servers []MServer

func Register(s MServer) {
	if s == nil {
		return
	}
	servers = append(servers, s)
}

func Start(countDown time.Duration) {

	// 启动server
	for _, server := range servers {
		go server.Start()
	}
	ctx := context.Background()
	defer func() {
		ctx.Done()
	}()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	// 宽容5s，应对elb和service对pod的调整（需要实测）
	<-time.After(5 * time.Second)
	// 向server发送shutdown信号
	for _, server := range servers {
		server.GracefulStop()
	}

	// 宽容countDown+2秒，等待各个server优雅停止，将已接收的请求执行完毕
	<-time.After(countDown + 2*time.Second)

}
