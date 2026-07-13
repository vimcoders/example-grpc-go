package echo

import (
	"context"
	"kube/generated/kubeapi"
)

type Handler struct {
}

func (h *Handler) Echo(ctx context.Context, req *kubeapi.HelloRequest) (*kubeapi.HelloResponse, error) {
	return &kubeapi.HelloResponse{Message: req.Message}, nil
}
