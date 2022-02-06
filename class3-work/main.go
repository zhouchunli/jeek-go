package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func HServer(srv *http.Server) error {
	addr := ":80"
	http.HandleFunc("/hello", HelloServer)
	srv.Addr = addr
	fmt.Println("http server start")
	err := srv.ListenAndServe()
	return err
}

// 增加一个 HTTP hanlder
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(ctx)
	srv := &http.Server{}
	g.Go(func() error {
		return HServer(srv)
	})
	g.Go(func() error {
		<-errCtx.Done()
		fmt.Println("http server done")
		return srv.Shutdown(errCtx)
	})
	g.Go(func() error {
		return fmt.Errorf("new error") // 模拟一个错误，导致整个group退出
	})
	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan)
	g.Go(func() error {
		for {
			select {
				case <-errCtx.Done():
					return errCtx.Err()
				case <-sChan:
					cancel()
			}
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Println("system error ", err)
	}
	fmt.Println("exist")
}