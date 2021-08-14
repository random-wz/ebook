package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ums/pkg/app/api"
)

var (
	sig     = make(chan os.Signal, 1)
	errChan = make(chan error, 1)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// 监控系统信号，确保程序有序退出
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go api.ListenAndServe(ctx, errChan)
	select {
	case err := <-errChan:
		log.Println("Error: http server | ", err.Error())
	case systemSig := <-sig:
		cancel()
		log.Println("exit ", systemSig.String())
	case <-ctx.Done():
		log.Println("ctx.Done")
	}
}
