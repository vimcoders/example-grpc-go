package gm

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.GMServiceServer
}

func (h *Handler) SendCommand(ctx context.Context, in *kubeapi.SendCommandRequest) (*kubeapi.SendCommandResponse, error) {
	slog.Info("GMServiceServer", "SendCommand", in)
	return &kubeapi.SendCommandResponse{}, nil
}
