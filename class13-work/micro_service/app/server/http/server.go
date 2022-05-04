package http

import (
	"context"
	"fmt"
	"micro_service/app/constant"
	"micro_service/app/deps"
	"micro_service/app/model/api"
	"micro_service/pkg/microerr"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	json "github.com/json-iterator/go"

	"github.com/julienschmidt/httprouter"
)

const (
	Get    = "GET"
	Post   = "POST"
	Put    = "PUT"
	Delete = "DELETE"
	Patch  = "PATCH"
	Option = "OPTIONS"
)

type Server struct {
	// 容量至少为1，保证GracefulStop发送消息时不会阻塞
	shutdown  chan struct{}
	srv       *http.Server
	countDown time.Duration
}

func NewServer(address string, countDown time.Duration) *Server {
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		//w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "application/json,charset=utf-8")
		//统一返回http的数据
		httpresp := api.HttpResponseBody{
			ApiStatus: int(microerr.ServerCommonError),
			SysStatus: int(microerr.Success),
			Info:      "Service Unavailable",
			Timestamp: int(time.Now().Unix()),
			Data:      api.EmptyRespData{},
		}
		_ = json.NewEncoder(w).Encode(httpresp)
		deps.Logger.Error(int(microerr.ServerCommonError),
			fmt.Sprintf("HTTP Server shutdown abnormally, Err: %v\n", v),
			string(debug.Stack()), "NA", constant.CommonHttpKey)
	}
	route(router)
	// app不支持301重定向，不让router自动产生重定向应答
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	return &Server{
		shutdown: make(chan struct{}, 1),
		srv: &http.Server{
			Addr:    address,
			Handler: router,
		},
		countDown: countDown,
	}
}

func (s *Server) Start() {
	go s.watchShutdown()
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("HTTP Server fail to start, Err: %v\n", err))
	}
}

func (s *Server) GracefulStop() {
	s.shutdown <- struct{}{}
}

func (s *Server) watchShutdown() {
	// 阻塞等待shutdown信号
	<-s.shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.countDown)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		deps.Logger.Error(int(microerr.ServerCommonError),
			fmt.Sprintf("HTTP Server Shutdown abnormally, Err: %v\n", err),
			"", "NA", constant.CommonHttpKey)
		if err = s.srv.Close(); err != nil { // 未测试 -- hhw
			deps.Logger.Error(int(microerr.ServerCommonError),
				fmt.Sprintf("HTTP Server Close abnormally, Err: %v\n", err),
				"", "NA", constant.CommonHttpKey)
		}
	}
}

func RunHTTPServer(address string) {
	httpListener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
	}

	router := httprouter.New()
	route(router)
	// app不支持301重定向，不让router自动产生重定向应答
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false

	_ = http.Serve(httpListener, router)
}
