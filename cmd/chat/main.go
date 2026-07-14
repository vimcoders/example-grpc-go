package main

import (
	"context"
	"example/app/chat"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.ChatService_ServiceDesc, &chat.Handler{})
	server.ListenAndServe(context.Background(), ":50051")
}
