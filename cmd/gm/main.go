package main

import (
	"context"
	"example/app/gm"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.GMService_ServiceDesc, &gm.Handler{})
	go func() {
		_ = s.ListenAndServe(ctx, ":50057")
		stop()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
