package main

import (
	"context"
	"kube/app/echo"
	"kube/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.HelloService_ServiceDesc, &echo.Handler{})
	server.ListenAndServe(context.Background(), ":50051")
}
