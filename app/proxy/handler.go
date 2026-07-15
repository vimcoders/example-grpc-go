package proxy

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.ProxyServiceServer
}

func (h *Handler) Login(ctx context.Context, in *kubeapi.LoginRequest) (*kubeapi.LoginResponse, error) {
	slog.Info("ProxyServiceServer", "Login", in)
	return &kubeapi.LoginResponse{}, nil
}
