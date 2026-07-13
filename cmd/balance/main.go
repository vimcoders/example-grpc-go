package main

import (
	"context"
	"kube/app/balance"

	"github.com/vimcoders/grpcx/generated/api"
)

func main() {
	server := balance.NewServer()
	server.RegisterService(&api.EchoService_ServiceDesc, ":50051")
	server.ListenAndServe(context.Background(), ":50052")
}
