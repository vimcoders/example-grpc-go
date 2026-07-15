package main

import (
	"context"
	"example/app/item"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.ItemService_ServiceDesc, &item.Handler{})
	server.ListenAndServe(context.Background(), ":50055")
}
