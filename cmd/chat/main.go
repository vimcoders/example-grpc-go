package main

import (
	"context"
	"example/app/chat"
	"example/generated/kubeapi"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.ChatService_ServiceDesc, &chat.Handler{})
	go func() {
		defer stop()
		_ = s.ListenAndServe(ctx, os.Getenv("Port"))
	}()
	slog.Info("running...")
	<-ctx.Done()
}
