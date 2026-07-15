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
	if err := server.RegisterService(&kubeapi.SocialService_ServiceDesc, "social:50052"); err != nil {
		panic(err)
	}
	if err := server.RegisterService(&kubeapi.ProxyService_ServiceDesc, "proxy:50053"); err != nil {
		panic(err)
	}
	if err := server.RegisterService(&kubeapi.ActivityService_ServiceDesc, "activity:50054"); err != nil {
		panic(err)
	}
	if err := server.RegisterService(&kubeapi.ItemService_ServiceDesc, "item:50055"); err != nil {
		panic(err)
	}
	if err := server.RegisterService(&kubeapi.MailService_ServiceDesc, "mail:50056"); err != nil {
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
