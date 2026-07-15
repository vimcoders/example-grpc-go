package chat

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.ChatServiceServer
}

func (h *Handler) Chat(ctx context.Context, req *kubeapi.ChatRequest) (*kubeapi.ChatResponse, error) {
	slog.Info("chat", "hello", req)
	return &kubeapi.ChatResponse{Message: req.Message}, nil
}
