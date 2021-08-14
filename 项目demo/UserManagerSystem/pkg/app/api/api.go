package api

import (
	"context"
	"net/http"
	"time"
	"ums/internal/app/api/routes"
)

func ListenAndServe(ctx context.Context, errChan chan error) {
	// TODO: 初始化数据库
	router := routes.InitRoutes()
	server := &http.Server{
		Addr:         ":10001",
		Handler:      router,
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 1000 * time.Second,
	}
	select {
	case errChan <- server.ListenAndServe():
		return
	case <-ctx.Done():
		return
	}
}
