package main

import (
	"context"
	"example/app/chat"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.ChatService_ServiceDesc, &chat.Handler{})
	go func() {
		_ = s.ListenAndServe(ctx, ":50051")
		stop()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
