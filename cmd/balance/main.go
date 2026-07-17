package main

import (
	"context"
	"example/app/balance"
	"example/generated/kubeapi"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := balance.NewServer(
		balance.WithRedisService(os.Getenv("RedisService")),
		balance.WithNatsService(os.Getenv("NatsService")),
	)
	for _, v := range []struct {
		endpoint string
		desc     grpc.ServiceDesc
	}{
		{os.Getenv("ChatService"), kubeapi.ChatService_ServiceDesc},
		{os.Getenv("SocialService"), kubeapi.SocialService_ServiceDesc},
		{os.Getenv("ProxyService"), kubeapi.ProxyService_ServiceDesc},
		{os.Getenv("ActivityService"), kubeapi.ActivityService_ServiceDesc},
		{os.Getenv("ItemService"), kubeapi.ItemService_ServiceDesc},
		{os.Getenv("MailService"), kubeapi.MailService_ServiceDesc},
		{os.Getenv("GMService"), kubeapi.GMService_ServiceDesc},
	} {
		if err := s.RegisterService(&v.desc, v.endpoint); err != nil {
			panic(err)
		}
	}
	go func() {
		defer stop()
		_ = s.ListenAndServe(ctx, os.Getenv("TCPPort"))
	}()
	go func() {
		defer stop()
		_ = s.ListenAndServeTLS(ctx, os.Getenv("TLSPort"))
	}()
	svr := &http.Server{
		Addr:           os.Getenv("HTTPPort"),
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: math.MaxInt16,
	}
	go func() {
		defer stop()
		_ = svr.ListenAndServe()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
