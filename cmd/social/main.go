package main

import (
	"context"
	"example/app/social"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.SocialService_ServiceDesc, &social.Handler{})
	go func() {
		_ = s.ListenAndServe(ctx, ":50052")
		stop()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
