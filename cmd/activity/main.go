package main

import (
	"context"
	"example/app/activity"
	"example/generated/kubeapi"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/vimcoders/grpcx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s := grpcx.NewServer()
	s.RegisterService(&kubeapi.ActivityService_ServiceDesc, &activity.Handler{})
	go func() {
		s.ListenAndServe(context.Background(), ":50054")
		stop()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
