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
	redisService, natsService := os.Getenv("RedisService"), os.Getenv("NatsService")
	slog.Info("NewServer", "redis", redisService, "nats", natsService)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	opt, err := balance.WithNatsService(natsService)
	if err != nil {
		slog.Error("WithNatsService", "error", err)
		return
	}
	opts := []balance.Option{opt, balance.WithRedisService(redisService)}
	s := balance.NewServer(opts...)
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
		slog.Info("RegisterService", "endpoint", v.endpoint)
		if err := s.RegisterService(v.desc, v.endpoint); err != nil {
			slog.Error("RegisterService", "endpoint", v.endpoint, "error", err)
			return
		}
	}
	if e := os.Getenv("TCPPort"); len(e) > 0 {
		slog.Info("TCP", "Port", e)
		go func() {
			defer stop()
			if err := s.ListenAndServe(ctx, e); err != nil {
				slog.Error("ListenAndServe", "stop", err)
			}
		}()
	}
	if e := os.Getenv("TLSPort"); len(e) > 0 {
		slog.Info("TLSPort", "Port", e)
		go func() {
			defer stop()
			if err := s.ListenAndServeTLS(ctx, e); err != nil {
				slog.Error("ListenAndServeTLS", "stop", err)
			}
		}()
	}
	if e := os.Getenv("HTTPPort"); len(e) > 0 {
		slog.Info("HTTPPort", "Port", e)
		svr := &http.Server{
			Addr:           e,
			Handler:        s,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    30 * time.Second,
			MaxHeaderBytes: math.MaxInt16,
		}
		go func() {
			defer stop()
			if err := svr.ListenAndServe(); err != nil {
				slog.Error("http ListenAndServe", "stop", err)
			}
			svr.Shutdown(ctx)
		}()
	}
	slog.Info("running...")
	<-ctx.Done()
}
