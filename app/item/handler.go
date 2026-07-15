package item

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.ItemServiceServer
}

func (h *Handler) Download(ctx context.Context, in *kubeapi.DownloadRequest) (*kubeapi.DownloadResponse, error) {
	slog.Info("ItemServiceServer", "Download", in)
	return &kubeapi.DownloadResponse{}, nil
}
