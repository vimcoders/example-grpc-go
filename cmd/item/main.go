package main

import (
	"context"
	"example/app/item"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.ItemService_ServiceDesc, &item.Handler{})
	go func() {
		defer stop()
		_ = s.ListenAndServe(ctx, ":50055")
	}()
	slog.Info("running...")
	<-ctx.Done()
}
