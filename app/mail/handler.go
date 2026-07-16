package mail

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.UnimplementedMailServiceServer
}

func (h *Handler) GetMail(ctx context.Context, in *kubeapi.GetMailRequest) (*kubeapi.GetMailResponse, error) {
	slog.Info("MailServiceServer", "GetMail", in)
	return &kubeapi.GetMailResponse{}, nil
}
