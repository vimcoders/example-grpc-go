package echo

import (
	"context"
	"kube/generated/api"
)

type Handler struct {
}

func (h *Handler) Echo(ctx context.Context, req *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: req.Message}, nil
}
