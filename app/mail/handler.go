package mail

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.MailServiceServer
}

func (h *Handler) GetMail(ctx context.Context, in *kubeapi.GetMailRequest) (*kubeapi.GetMailResponse, error) {
	slog.Info("MailServiceServer", "GetMail", in)
	return &kubeapi.GetMailResponse{}, nil
}
