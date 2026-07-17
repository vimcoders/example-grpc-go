package main

import (
	"context"
	"example/app/mail"
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
	s.RegisterService(&kubeapi.MailService_ServiceDesc, &mail.Handler{})
	go func() {
		defer stop()
		_ = s.ListenAndServe(ctx, os.Getenv("Port"))
	}()
	slog.Info("running...")
	<-ctx.Done()
}
