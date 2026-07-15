package main

import (
	"context"
	"example/app/balance"
	"example/generated/kubeapi"
	"math"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	server := balance.NewServer()
	if err := server.RegisterService(&kubeapi.ChatService_ServiceDesc, "chat:50051"); err != nil {
		panic(err)
	}
	go func() {
		_ = server.ListenAndServe(ctx, ":26888")
	}()
	svr := &http.Server{
		Addr:           ":36888",
		Handler:        server,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: math.MaxInt16,
	}
	go func() {
		_ = svr.ListenAndServe()
	}()
	<-ctx.Done()
}
