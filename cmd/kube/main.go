package main

import (
	"context"
	"example/app/kube"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.KubeService_ServiceDesc, &kube.Handler{})
	server.ListenAndServe(context.Background(), ":50051")
}
