package grpc

import (
	"fmt"
	"micro_service/app/constant"
	controller2 "micro_service/app/controller"
	"micro_service/app/deps"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/grpc/keepalive"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
)

type Server struct {
	// 容量至少为1，保证GracefulStop发送消息时不会阻塞
	shutdown chan struct{}
	srv      *grpc.Server
	// GracefulStop的宽容时间
	countDown time.Duration
	addr      string
}

func NewServer(address string, countDown time.Duration) *Server {

	// 定义panicHandler
	panicHandler := func(p interface{}) (err error) {
		deps.Logger.Error(550, "NA", string(debug.Stack()), "NA", "NA")
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(panicHandler),
	}

	// keepalive配置
	kaep := keepalive.EnforcementPolicy{
		// 需要小于客户端 keepalive.ClientParameters.Time
		MinTime:             10 * time.Second, // If a client pings more than once every 10 seconds, terminate the connection
		PermitWithoutStream: true,             // Allow pings even when there are no active streams
	}
	kasp := keepalive.ServerParameters{
		// 最大空闲连接时长
		MaxConnectionIdle: time.Minute, // If a client is idle for 1 minute, send a GOAWAY
		// 最大连接时长 暂不设置 默认永久
		//MaxConnectionAge: 60 * time.Second, // If any connection is alive for more than 60 seconds, send a GOAWAY
		// 强制关闭前，等待通信(阻塞/处理)中的请求完成
		MaxConnectionAgeGrace: 5 * time.Second, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		// 服务端暂不配置主动ping规则
		//Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		//Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}

	baseServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(opts...),
		)),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	)

	// 注册request handler
	// HealthService 健康检查
	healthpb.RegisterHealthServer(baseServer, &controller2.HealthService{})

	return &Server{
		shutdown:  make(chan struct{}, 1),
		srv:       baseServer,
		countDown: countDown,
		addr:      address,
	}
}

func (s *Server) Start() {
	go s.watchShutdown() // 启动监控
	grpcListener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(fmt.Sprintf("GRPC Server fail to listen, Addr: %s, Err: %v\n", s.addr, err))
	}
	reflection.Register(s.srv)
	err = s.srv.Serve(grpcListener)
	if err != nil && err != grpc.ErrServerStopped {
		panic(fmt.Sprintf("GRPC Server fail to start, Addr: %s, Err: %v\n", s.addr, err))
	}
}

// 向shutdown通道发送信号
func (s *Server) GracefulStop() {
	s.shutdown <- struct{}{}
}

func (s *Server) watchShutdown() {
	// 阻塞等待shutdown信号
	<-s.shutdown
	// 容量为1 ，防止通道阻塞，goroutine泄露
	cancel := make(chan struct{}, 1)
	go func() {
		// 超过宽容时间，强制关闭
		t := time.NewTimer(s.countDown)
		select {
		case <-t.C:
			s.srv.Stop() // 未测试 -- hhw
		case <-cancel:
			t.Stop()
		}
	}()
	s.srv.GracefulStop()
	// 成功关闭，取消定时器
	cancel <- struct{}{}
}
