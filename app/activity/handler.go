package activity

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.UnimplementedActivityServiceServer
}

func (h *Handler) GetActivity(ctx context.Context, in *kubeapi.GetActivityRequest) (*kubeapi.GetActivityResponse, error) {
	slog.Info("ActivityServiceServer", "GetActivity", in)
	return &kubeapi.GetActivityResponse{}, nil
}
