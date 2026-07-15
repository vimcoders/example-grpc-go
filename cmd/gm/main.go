package main

import (
	"context"
	"example/app/gm"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.GMService_ServiceDesc, &gm.Handler{})
	server.ListenAndServe(context.Background(), ":50057")
}
