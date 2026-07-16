package main

import (
	"context"
	"example/app/proxy"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.ProxyService_ServiceDesc, &proxy.Handler{})
	go func() {
		_ = s.ListenAndServe(ctx, ":50053")
		stop()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
