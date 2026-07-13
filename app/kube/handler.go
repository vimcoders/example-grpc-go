package kube

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.KubeServiceServer
}

func (h *Handler) KubeEcho(ctx context.Context, req *kubeapi.KubeRequest) (*kubeapi.KubeResponse, error) {
	slog.Info("kube", "hello", req)
	return &kubeapi.KubeResponse{Message: req.Message}, nil
}
