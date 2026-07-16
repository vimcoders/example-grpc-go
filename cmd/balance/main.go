package main

import (
	"context"
	"example/app/balance"
	"log/slog"
	"math"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	server := balance.NewServer()
	go func() {
		defer stop()
		_ = server.ListenAndServe(ctx, ":26888")
	}()
	go func() {
		defer stop()
		_ = server.ListenAndServeTLS(ctx, ":56888")
	}()
	svr := &http.Server{
		Addr:           ":36888",
		Handler:        server,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: math.MaxInt16,
	}
	go func() {
		defer stop()
		_ = svr.ListenAndServe()
	}()
	slog.Info("running...")
	<-ctx.Done()
}
