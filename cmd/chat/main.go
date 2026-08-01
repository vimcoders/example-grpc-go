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
	port := os.Getenv("Port")
	go func() {
		defer stop()
		if err := s.ListenAndServe(ctx, port); err != nil {
			slog.Error("ListenAndServe", "Port", port)
		}
	}()
	slog.Info("running...", "Port", port)
	<-ctx.Done()
}
