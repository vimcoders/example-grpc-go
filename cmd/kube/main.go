package main

import (
	"context"
	"kube/app/kube"
	"kube/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.KubeService_ServiceDesc, &kube.Handler{})
	server.ListenAndServe(context.Background(), ":50051")
}
