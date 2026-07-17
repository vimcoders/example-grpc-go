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
	s := balance.NewServer()
	go func() {
		defer stop()
		_ = s.ListenAndServe(ctx, ":26888")
	}()
	go func() {
		defer stop()
		_ = s.ListenAndServeTLS(ctx, ":56888")
	}()
	svr := &http.Server{
		Addr:           ":16888",
		Handler:        s,
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
