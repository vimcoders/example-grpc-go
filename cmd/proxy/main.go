package main

import (
	"context"
	"example/app/proxy"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.ProxyService_ServiceDesc, &proxy.Handler{})
	server.ListenAndServe(context.Background(), ":50053")
}
