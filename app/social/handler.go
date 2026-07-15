package social

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
)

type Handler struct {
	kubeapi.SocialServiceServer
}

func (h *Handler) GetFriends(ctx context.Context, in *kubeapi.GetFriendsRequest) (*kubeapi.GetFriendsResponse, error) {
	slog.Info("SocialServiceServer", "GetFriends", in)
	return &kubeapi.GetFriendsResponse{}, nil
}
